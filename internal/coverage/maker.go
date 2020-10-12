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
	"reflect"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/MattWindsor91/act-tester/internal/stage/lifter"

	"github.com/MattWindsor91/act-tester/internal/stage/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/observing"

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

	// rng is the random number generator used to seed the various profile makers.
	rng *rand.Rand

	// errw is the writer to which we send stderr from standalone coverage makers.
	errw io.Writer

	// observers contains the observers being used by the maker.  Each is accessed in at most one thread at a time.
	observers []Observer
}

// NewMaker constructs a new coverage testbed maker.
func NewMaker(outDir string, profiles map[string]Profile, opts ...Option) (*Maker, error) {
	m := &Maker{outDir: outDir, profiles: profiles, rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
	if err := Options(opts...)(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Maker) Run(ctx context.Context) error {
	buckets := m.qs.Buckets()
	if buckets == nil {
		return errors.New("bucket calculation failed")
	}
	return m.runProfiles(ctx, buckets)
}

func (m *Maker) runProfiles(ctx context.Context, buckets map[string]int) error {
	obsChs := make(map[string]<-chan RunMessage, len(m.profiles))

	// TODO(@MattWindsor91): worker queue
	eg, ectx := errgroup.WithContext(ctx)
	for pname, p := range m.profiles {
		obsCh := make(chan RunMessage)

		pm, err := m.makeProfileMaker(pname, p, buckets, obsCh)
		if err != nil {
			return err
		}

		eg.Go(func() error {
			err := pm.run(ectx)
			close(obsCh)
			return err
		})

		obsChs[pname] = obsCh
	}
	eg.Go(func() error {
		return m.fanInObservations(ectx, obsChs)
	})
	return eg.Wait()
}

func (m *Maker) fanInObservations(ectx context.Context, obsChs map[string]<-chan RunMessage) error {
	cs := make([]reflect.SelectCase, len(obsChs)+1)
	cs[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ectx.Done())}
	i := 1
	for _, och := range obsChs {
		cs[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(och)}
		i++
	}

	// TODO(@MattWindsor91): consider generalising/replicating this fan-in pattern
	want := len(m.profiles)
	for {
		chosen, recv, recvOK := reflect.Select(cs)
		if chosen == 0 {
			// Drain every channel.
			for _, och := range obsChs {
				for range och {
				}
			}
			return ectx.Err()
		}
		if !recvOK {
			continue
		}
		rm := recv.Interface().(RunMessage)
		OnCoverageRun(rm, m.observers...)
		if rm.Kind == observing.BatchEnd {
			want--
			if want == 0 {
				return nil
			}
		}
	}
}

func (m *Maker) makeProfileMaker(pname string, p Profile, buckets map[string]int, obsCh chan<- RunMessage) (*profileMaker, error) {
	runner, err := m.makeRunner(p)
	if err != nil {
		return nil, err
	}

	pm := &profileMaker{
		name:    pname,
		dir:     filepath.Join(m.outDir, pname),
		profile: p,
		buckets: buckets,
		total:   m.qs.Count,
		runner:  runner,
		obsCh:   obsCh,
		rng:     rand.New(rand.NewSource(m.rng.Int63())),
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
		ErrW:       m.errw,
	}, nil
}
