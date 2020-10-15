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
	buckets []Bucket
	inputs  []string
	total   int
	runner  Runner
	rng     *rand.Rand
	jobCh   chan<- workerJob

	// index of current step in the profile run; we use 32-bit to make it possible to serve this as an act-fuzz seed.
	nrun int32
}

func (p *profileMaker) run(ctx context.Context) error {
	// Observations are handled in the worker.
	p.nrun = 0
	for _, bucket := range p.buckets {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := p.runBucket(ctx, bucket); err != nil {
			return err
		}
	}
	return nil
}

func (p *profileMaker) runBucket(ctx context.Context, b Bucket) error {
	for i := 0; i < b.Size; i++ {
		p.nrun++

		rc, err := p.makeRunContext(b.Name, i)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case p.jobCh <- p.makeJob(rc):
		}
	}
	return nil
}

func (p *profileMaker) makeJob(rc RunContext) workerJob {
	return workerJob{
		pname:   p.name,
		profile: p.profile,
		nrun:    int(p.nrun),
		rc:      rc,
		r:       p.runner,
	}
}

func (p *profileMaker) makeRunContext(suffix string, i int) (RunContext, error) {
	in, err := p.randomInput()
	return RunContext{
		Seed:        p.nrun,
		BucketDir:   filepath.Join(p.dir, suffix),
		NumInBucket: i,
		Input:       in,
	}, err
}

func (p *profileMaker) mkdirs() error {
	for _, b := range p.buckets {
		bdir := p.bucketDir(b.Name)
		if err := yos.MakeDir(bdir); err != nil {
			return fmt.Errorf("preparing directory for profile %q bucket %q: %w", p.name, b, err)
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
