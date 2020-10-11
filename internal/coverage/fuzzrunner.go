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
	"path/filepath"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	litmus2 "github.com/MattWindsor91/act-tester/internal/model/litmus"
	"github.com/MattWindsor91/act-tester/internal/model/service"

	"github.com/MattWindsor91/act-tester/internal/stage/lifter"

	"github.com/MattWindsor91/act-tester/internal/model/job"
	fuzzer2 "github.com/MattWindsor91/act-tester/internal/model/service/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/stage/fuzzer"
)

var (
	// TODO(@MattWindsor91): this feels like it's duplicating a lot of logic elsewhere in the tester.

	// ErrNoArch occurs when we try to run a lifted fuzzer for coverage and haven't a target architecture.
	ErrNoArch = errors.New("this runner needs an architecture set, but got none")
	// ErrNoFuzzer occurs when we try to run a lifted fuzzer for coverage but the fuzzer is nil.
	ErrNoFuzzer = errors.New("this runner needs a fuzzer, but got none")
	// ErrNoLifter occurs when we try to run a lifted fuzzer for coverage but the lifter is nil.
	ErrNoLifter = errors.New("this runner needs a lifter, but got none")
	// ErrNoInput occurs when we try to run a mutating fuzzer for coverage and haven't any input to feed it.
	ErrNoInput = errors.New("this runner needs input testcases, but got none")
)

// FuzzRunner is a coverage runner that uses the act fuzzer.
type FuzzRunner struct {
	// Fuzzer is the fuzzer this fuzz runner uses.
	Fuzzer fuzzer.SingleFuzzer
	// Lifter is the lifter this lifter uses.
	Lifter lifter.SingleLifter
	// Config is the configuration to pass to the fuzz runner.
	Config *fuzzer2.Configuration

	// Arch is the architecture that the lifting process should target.
	Arch id.ID
	// Backend can point to the backend information for the lifter.
	Backend *service.Backend
	// ErrW can point to a writer that should receive stderr from any external programs.
	ErrW io.Writer
}

func (f *FuzzRunner) Run(ctx context.Context, rc RunContext) error {
	// TODO(@MattWindsor91): this should probably be a smart constructor
	if err := f.check(); err != nil {
		return fmt.Errorf("fuzz runner internal checks failed: %w", err)
	}

	litmus, err := f.fuzz(ctx, rc)
	if err != nil {
		return fmt.Errorf("while fuzzing (%q -> %q): %w", rc.inputPathOrEmpty(), rc.OutLitmus(), err)
	}
	return f.lift(ctx, rc, litmus)
}

func (f *FuzzRunner) check() error {
	if f.Fuzzer == nil {
		return ErrNoFuzzer
	}
	if f.Lifter == nil {
		return ErrNoLifter
	}
	if f.Arch.IsEmpty() {
		return ErrNoArch
	}
	return nil
}

func (f *FuzzRunner) fuzz(ctx context.Context, rc RunContext) (string, error) {
	j, err := f.fuzzJob(rc)
	if err != nil {
		return "", fmt.Errorf("can't make fuzzer job: %w", err)
	}
	return j.OutLitmus, f.Fuzzer.Fuzz(ctx, j)
}

func (f *FuzzRunner) fuzzJob(rc RunContext) (job.Fuzzer, error) {
	in, err := rc.inputPath()
	if err != nil {
		return job.Fuzzer{}, err
	}

	return job.Fuzzer{
		Seed:      rc.Seed,
		In:        filepath.ToSlash(in),
		OutLitmus: rc.OutLitmus(),
		Config:    f.Config,
	}, nil
}

func (f *FuzzRunner) lift(ctx context.Context, rc RunContext, litmus string) error {
	// TODO(@MattWindsor91): do something with the recipe
	_, err := f.Lifter.Lift(ctx, f.liftJob(litmus, rc), f.ErrW)
	return err
}

func (f *FuzzRunner) liftJob(litmus string, rc RunContext) job.Lifter {
	return job.Lifter{
		Backend: f.Backend,
		Arch:    f.Arch,
		In:      *litmus2.New(filepath.ToSlash(litmus)),
		OutDir:  rc.LiftOutDir(),
	}
}
