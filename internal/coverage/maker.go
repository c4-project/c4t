// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"path/filepath"

	"github.com/MattWindsor91/act-tester/internal/helper/srvrun"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"
	"github.com/MattWindsor91/act-tester/internal/observing"

	"github.com/MattWindsor91/act-tester/internal/stage/lifter"

	"github.com/MattWindsor91/act-tester/internal/stage/fuzzer"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrNeedBackend occurs when we try to instantiate a runner for a known-fuzzer profile without a backend.
	ErrNeedBackend = errors.New("need backend information for this profile")

	// ErrNeedRunInfo occurs when we try to instantiate a runner for a standalone profile without run information.
	ErrNeedRunInfo = errors.New("need run information for this profile")

	// ErrUnsupportedProfileKind occurs when we try to instantiate a runner for an unsupported profile type.
	ErrUnsupportedProfileKind = errors.New("this profile kind can't be run yet")
)

// Maker contains state used by the coverage testbed maker.
type Maker struct {
	// outDir is the name of the output directory.
	outDir string

	// profiles contains the map of profiles available to the coverage testbed maker.
	profiles map[string]Profile

	// TODO(@MattWindsor91): Add multiple fuzzers and lifters

	// fuzz tells the maker how to run its internal fuzzer.
	fuzz fuzzer.SingleFuzzer

	// lift tells the maker how to run its internal lifter.
	lift lifter.SingleLifter

	// sdump tells the maker how to dump statistics.
	sdump litmus.StatDumper

	// qs is the calculated quantity set for the coverage testbed maker.
	qs QuantitySet

	// inputs contains the filepaths to each input subject to use for fuzzing profiles that need them.
	inputs []string

	// errw is the writer to which we send stderr from standalone coverage makers.
	errw io.Writer

	// observers contains the observers being used by the maker.  Each is accessed in at most one thread at a time.
	observers []Observer

	// fanIn handles fan-in for observers across other threads.
	fanIn *observing.FanIn
}

const (
	// DefaultCount is the default value for the Count quantity.
	DefaultCount = 1000
	// DefaultNWorkers is the default value for the NWorkers quantity.
	DefaultNWorkers = 10
)

// NewMaker constructs a new coverage testbed maker.
func NewMaker(outDir string, profiles map[string]Profile, opts ...Option) (*Maker, error) {
	m := &Maker{
		outDir:   outDir,
		profiles: profiles,
		qs: QuantitySet{
			Count:     DefaultCount,
			Divisions: nil,
			NWorkers:  DefaultNWorkers,
		},
	}
	if err := Options(opts...)(m); err != nil {
		return nil, err
	}
	m.fanIn = observing.NewFanIn(func(_ int, input interface{}) error {
		OnCoverageRun(input.(RunMessage), m.observers...)
		return nil
	}, m.qs.NWorkers)
	return m, nil
}

func (m *Maker) Run(ctx context.Context) error {
	buckets := m.qs.Buckets()
	if buckets == nil {
		return errors.New("bucket calculation failed")
	}
	return m.runProfiles(ctx, buckets)
}

func (m *Maker) runProfiles(ctx context.Context, buckets []Bucket) error {
	// TODO(@MattWindsor91): I'm not sure what the correct value here should be.
	jobCh := make(chan workerJob, m.qs.NWorkers)

	// TODO(@MattWindsor91): worker queue
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return m.feeder(ctx, buckets, jobCh)
	})
	for i := 0; i < m.qs.NWorkers; i++ {
		// TODO(@MattWindsor91): sensible buffer size?
		obsCh := make(chan RunMessage, m.qs.NWorkers)
		eg.Go(func() error {
			return m.worker(ctx, obsCh, jobCh)
		})
		m.fanIn.Add(obsCh)
	}
	eg.Go(func() error {
		return m.fanIn.Run(ectx)
	})
	return eg.Wait()
}

type workerJob struct {
	pname   string
	profile Profile
	nrun    int
	rc      RunContext
	r       Runner
}

func (m *Maker) feeder(ctx context.Context, buckets []Bucket, jobCh chan<- workerJob) error {
	// This errgroup nesting mainly exists so we only need one job channel.
	// TODO(@MattWindsor91): can we do better here?
	defer close(jobCh)

	eg, ectx := errgroup.WithContext(ctx)
	for pname, p := range m.profiles {
		pm, err := m.makeProfileMaker(pname, p, buckets, jobCh)
		if err != nil {
			return err
		}
		eg.Go(func() error {
			return pm.run(ectx)
		})
	}
	return eg.Wait()
}

func (m *Maker) worker(ctx context.Context, obsCh chan<- RunMessage, jobCh <-chan workerJob) error {
	defer close(obsCh)

	for {
		select {
		case <-ctx.Done():
			for range jobCh {
			}
			return ctx.Err()
		case j, ok := <-jobCh:
			if !ok {
				return nil
			}
			if err := m.workerJob(ctx, obsCh, j); err != nil {
				return err
			}
		}
	}
}

func (m *Maker) workerJob(ctx context.Context, obsCh chan<- RunMessage, j workerJob) error {
	// nrun is 1-indexed
	if j.nrun == 1 {
		obsCh <- RunStart(j.pname, j.profile, m.qs.Count)
	}
	obsCh <- RunStep(j.pname, j.nrun, j.rc)
	if err := j.r.Run(ctx, j.rc); err != nil {
		return err
	}
	if j.nrun == m.qs.Count {
		obsCh <- RunEnd(j.pname)
	}
	return nil
}

func (m *Maker) makeProfileMaker(pname string, p Profile, buckets []Bucket, jobCh chan<- workerJob) (*profileMaker, error) {
	runner, err := m.makeRunner(p)
	if err != nil {
		return nil, err
	}

	// The idea here is to have something that is technically deterministic, but tours the input space in a seemingly
	// random order.  Why?  Because for inputs like Memalloy, input number tends to correlate with input complexity,
	// and we don't want to give all the simple inputs to the first runs and the complex ones to the later ones.
	rng := rand.New(rand.NewSource(0))

	pm := &profileMaker{
		name:    pname,
		dir:     filepath.Join(m.outDir, pname),
		profile: p,
		buckets: buckets,
		total:   m.qs.Count,
		runner:  runner,
		jobCh:   jobCh,
		rng:     rng,
		inputs:  m.inputs,
	}

	if err := pm.mkdirs(); err != nil {
		return nil, fmt.Errorf("preparing directories for profile %q: %w", pname, err)
	}
	return pm, nil
}

func (m *Maker) makeRunner(p Profile) (Runner, error) {
	// this mostly used only for testing
	if p.Runner != nil {
		return p.Runner, nil
	}

	switch p.Kind {
	case Standalone:
		return m.standaloneRunner(p)
	case Known:
		return m.knownRunner(p)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProfileKind, p.Kind)
	}
}

func (m *Maker) standaloneRunner(p Profile) (*StandaloneRunner, error) {
	if p.Run == nil {
		return nil, ErrNeedRunInfo
	}
	return &StandaloneRunner{run: *p.Run, errw: m.errw}, nil
}

func (m *Maker) knownRunner(p Profile) (Runner, error) {
	if p.Backend == nil {
		return nil, ErrNeedBackend
	}
	return &FuzzRunner{
		Fuzzer:     m.fuzz,
		Lifter:     m.lift,
		StatDumper: m.sdump,
		Config:     p.Fuzz,
		Arch:       p.Arch,
		Backend:    p.Backend,
		// TODO(@MattWindsor91): push this up somehow
		Runner: srvrun.NewExecRunner(srvrun.StderrTo(m.errw)),
	}, nil
}
