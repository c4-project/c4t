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

	"github.com/MattWindsor91/act-tester/internal/observing"

	"golang.org/x/sync/errgroup"

	"github.com/1set/gut/yos"
)

// Maker contains state used by the coverage testbed maker.
type Maker struct {
	// outDir is the name of the output directory.
	outDir string

	// profiles contains the map of profiles available to the coverage testbed maker.
	profiles map[string]Profile

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

	if err := m.prepare(buckets); err != nil {
		return err
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
	})
	return eg.Wait()
}

func (m *Maker) makeProfileMaker(pname string, p Profile, buckets map[string]int, obsCh chan<- RunMessage) (*profileMaker, error) {
	runner, err := p.runner(m.errw)
	if err != nil {
		return nil, err
	}

	return &profileMaker{
		name:    pname,
		dir:     filepath.Join(m.outDir, pname),
		profile: p,
		buckets: buckets,
		total:   m.qs.Count,
		runner:  runner,
		obsCh:   obsCh,
		rng:     rand.New(rand.NewSource(m.rng.Int63())),
	}, nil
}

func (m *Maker) prepare(buckets map[string]int) error {
	for pname := range m.profiles {
		for suffix := range buckets {
			if err := yos.MakeDir(m.bucketDir(pname, suffix)); err != nil {
				return fmt.Errorf("preparing directory for profile %q bucket %q: %w", pname, suffix, err)
			}
		}
	}
	return nil
}

func (m *Maker) bucketDir(pname string, suffix string) string {
	return filepath.Join(m.outDir, pname, suffix)
}
