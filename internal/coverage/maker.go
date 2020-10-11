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
	"time"

	"github.com/MattWindsor91/act-tester/internal/stage/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/observing"

	"golang.org/x/sync/errgroup"
)

// Maker contains state used by the coverage testbed maker.
type Maker struct {
	// outDir is the name of the output directory.
	outDir string

	// profiles contains the map of profiles available to the coverage testbed maker.
	profiles map[string]Profile

	// TODO(@MattWindsor91): Add multiple fuzzers

	// fuzz tells the maker how to run act-fuzz.
	fuzz fuzzer.SingleFuzzer

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
	obs := make(chan RunMessage)

	// TODO(@MattWindsor91): worker queue
	eg, ectx := errgroup.WithContext(ctx)
	for pname, p := range m.profiles {
		pm, err := m.makeProfileMaker(pname, p, buckets, obs)
		if err != nil {
			return err
		}

		eg.Go(func() error {
			return pm.run(ectx)
		})
	}
	eg.Go(func() error {
		return m.fanInObservations(ectx, obs)
	})
	return eg.Wait()
}

func (m *Maker) fanInObservations(ectx context.Context, obs <-chan RunMessage) error {
	// TODO(@MattWindsor91): consider generalising/replicating this fan-in pattern
	want := len(m.profiles)
	for {
		select {
		case <-ectx.Done():
			return ectx.Err()
		case rm := <-obs:
			OnCoverageRun(rm, m.observers...)
			if rm.Kind == observing.BatchEnd {
				want--
				if want == 0 {
					return nil
				}
			}
		}
	}
}

func (m *Maker) makeProfileMaker(pname string, p Profile, buckets map[string]int, obsCh chan<- RunMessage) (*profileMaker, error) {
	runner, err := p.runner(m.errw)
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
	}

	if err := pm.mkdirs(); err != nil {
		return nil, fmt.Errorf("preparing directories for profile %q: %w", pname, err)
	}
	return pm, nil
}
