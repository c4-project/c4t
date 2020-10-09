// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"context"
	"math/rand"
	"path/filepath"
)

// profileMaker governs the making of a testbed under one particular coverage profile.
type profileMaker struct {
	name    string
	dir     string
	profile Profile
	buckets map[string]int
	total   int
	runner  Runner
	rng     *rand.Rand
	obsCh   chan<- RunMessage

	// index of current step in the profile run
	nstep int
}

func (p *profileMaker) run(ctx context.Context) error {
	p.obsCh <- RunStart(p.name, p.profile, p.total)

	p.nstep = 0
	for suffix, bsize := range p.buckets {
		if err := p.runBucket(ctx, bsize, suffix); err != nil {
			return err
		}
	}

	p.obsCh <- RunEnd(p.name)
	return nil
}

func (p *profileMaker) runBucket(ctx context.Context, bsize int, suffix string) error {
	for i := 0; i < bsize; i++ {
		p.nstep++

		rc := p.makeRunContext(suffix, i)
		p.obsCh <- RunStep(p.nstep, rc)
		if err := p.runner.Run(ctx, rc); err != nil {
			return err
		}
	}
	return nil
}

func (p *profileMaker) makeRunContext(suffix string, i int) RunnerContext {
	return RunnerContext{
		// Using 32-bit seed for compatibility with things like act-fuzz.
		Seed:        p.rng.Int31(),
		BucketDir:   filepath.Join(p.dir, suffix),
		NumInBucket: i,
		Input:       nil,
	}
}
