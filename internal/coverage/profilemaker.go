// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"context"
	"fmt"
	"math/rand"
	"path/filepath"

	"github.com/1set/gut/yos"
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

func (p *profileMaker) mkdirs() error {
	for suffix := range p.buckets {
		bdir := p.bucketDir(suffix)
		if err := yos.MakeDir(bdir); err != nil {
			return fmt.Errorf("preparing directory for profile %q bucket %q: %w", p.name, suffix, err)
		}
	}
	return nil
}

func (p *profileMaker) bucketDir(suffix string) string {
	return filepath.Join(p.dir, suffix)
}
