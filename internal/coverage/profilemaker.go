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

	"github.com/MattWindsor91/act-tester/internal/model/litmus"
	"github.com/MattWindsor91/act-tester/internal/subject"

	"github.com/1set/gut/yos"
)

// profileMaker governs the making of a testbed under one particular coverage profile.
type profileMaker struct {
	name    string
	dir     string
	profile Profile
	buckets map[string]int
	inputs  []string
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
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := p.runBucket(ctx, bsize, suffix); err != nil {
			return err
		}
	}

	p.obsCh <- RunEnd(p.name)
	return nil
}

func (p *profileMaker) runBucket(ctx context.Context, bsize int, suffix string) error {
	for i := 0; i < bsize; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		p.nstep++

		rc, err := p.makeRunContext(suffix, i)
		if err != nil {
			return err
		}
		p.obsCh <- RunStep(p.name, p.nstep, rc)
		if err := p.runner.Run(ctx, rc); err != nil {
			return err
		}
	}
	return nil
}

func (p *profileMaker) makeRunContext(suffix string, i int) (RunContext, error) {
	in, err := p.randomInput()
	return RunContext{
		// Using 32-bit seed for compatibility with things like act-fuzz.
		Seed:        p.rng.Int31(),
		BucketDir:   filepath.Join(p.dir, suffix),
		NumInBucket: i,
		Input:       in,
	}, err
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

func (p *profileMaker) randomInput() (*subject.Subject, error) {
	nin := len(p.inputs)
	// TODO(@MattWindsor91): probe the subject properly?
	if nin == 0 {
		return nil, nil
	}
	input := p.inputs[p.rng.Intn(nin)]
	return subject.New(litmus.New(filepath.ToSlash(input)))
}
