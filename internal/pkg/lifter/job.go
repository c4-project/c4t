// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package lifter

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// Job is the type of per-architecture lifter jobs.
type Job struct {
	// Arch is the architecture for which this job is responsible.
	Arch model.ID

	// Backend is the ID of the backend that this job will use.
	Backend model.ID

	// Maker is the harness maker for this job.
	Maker HarnessMaker

	// OutDir is the root output directory for this lifter job.
	OutDir string

	// Corpus is the existing corpus that we are trying to lift.
	Corpus corpus.Corpus

	// Rng is the random number generator to use for fuzz seeds.
	Rng *rand.Rand

	// ResCh is the channel onto which each fuzzed subject should be sent.
	ResCh chan<- builder.Request
}

// Lift performs this lifting job.
func (j *Job) Lift(ctx context.Context) error {
	if err := j.check(); err != nil {
		return err
	}

	return j.Corpus.Par(ctx, 20, func(ctx context.Context, s subject.Named) error {
		return j.liftSubject(ctx, &s)
	})
}

// check does some basic checking on the Job before starting to run it.
func (j *Job) check() error {
	if j.Backend.IsEmpty() {
		return ErrNoBackend
	}
	if j.Corpus == nil {
		return corpus.ErrNone
	}
	if j.Maker == nil {
		return ErrMakerNil
	}
	return nil
}

func (j *Job) liftSubject(ctx context.Context, s *subject.Named) error {
	// TODO(@MattWindsor91): bring this in line with the other stages' pathsets
	dir, derr := buildAndMkDir(j.OutDir, s.Name)
	if derr != nil {
		return fmt.Errorf("when making subject dir: %w", derr)
	}

	path, perr := s.BestLitmus()
	if perr != nil {
		return perr
	}

	spec := model.HarnessSpec{
		Backend: j.Backend,
		Arch:    j.Arch,
		InFile:  path,
		OutDir:  dir,
	}

	files, err := j.Maker.MakeHarness(ctx, spec)
	if err != nil {
		return fmt.Errorf("when making harness for %s (arch %s): %w", s.Name, j.Arch.String(), err)
	}

	return j.makeBuilderReq(s, dir, files).SendTo(ctx, j.ResCh)
}

func (j *Job) makeBuilderReq(s *subject.Named, dir string, files []string) builder.Request {
	return builder.Request{
		Name: s.Name,
		Req: builder.Harness{
			Arch:    j.Arch,
			Harness: subject.Harness{Dir: dir, Files: files},
		},
	}
}
