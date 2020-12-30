// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/c4-project/c4t/internal/model/service"
	"github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/1set/gut/yos"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/model/litmus"

	fuzzer2 "github.com/c4-project/c4t/internal/model/service/fuzzer"
	"github.com/c4-project/c4t/internal/stage/fuzzer"
)

var (
	// TODO(@MattWindsor91): this feels like it's duplicating a lot of logic elsewhere in the tester.
	// TODO(@MattWindsor91): this actually feels like a good starting point for a single-file runner for bisection etc.

	// ErrNoArch occurs when we try to run a lifted fuzzer for coverage and haven't a target architecture.
	ErrNoArch = errors.New("this runner needs an architecture set, but got none")
	// ErrNoFuzzer occurs when we try to run a lifted fuzzer for coverage but the fuzzer is nil.
	ErrNoFuzzer = errors.New("this runner needs a fuzzer, but got none")
	// ErrNoLifter occurs when we try to run a lifted fuzzer for coverage but the lifter is nil.
	ErrNoLifter = errors.New("this runner needs a lifter, but got none")
	// ErrNoStatDumper occurs when we try to run a lifted fuzzer for coverage but the statistic dumper is nil.
	ErrNoStatDumper = errors.New("this runner needs a lifter, but got none")
	// ErrNoInput occurs when we try to run a mutating fuzzer for coverage and haven't any input to feed it.
	ErrNoInput = errors.New("this runner needs input testcases, but got none")
)

// FuzzRunner is a coverage runner that uses the act fuzzer.
type FuzzRunner struct {
	// TODO(@MattWindsor91): this overlaps with Env in director.
	// Fuzzer is the fuzzer this fuzz runner uses.
	Fuzzer fuzzer.SingleFuzzer
	// Lifter is the lifter this fuzz runner uses.
	Lifter backend.SingleLifter
	// StatDumper is the statistics dumper this fuzz runner uses between fuzzing and lifting.
	StatDumper litmus.StatDumper
	// Config is the configuration to pass to the fuzz runner.
	Config *fuzzer2.Configuration

	// Arch is the architecture that the lifting process should target.
	Arch id.ID
	// Runner should be the service runner to use when invoking the lifter.
	Runner service.Runner
}

func (f *FuzzRunner) Run(ctx context.Context, rc RunContext) error {
	// TODO(@MattWindsor91): this should probably be a smart constructor
	if err := f.check(); err != nil {
		return fmt.Errorf("fuzz runner internal checks failed: %w", err)
	}

	fuzzed, err := f.fuzz(ctx, rc)
	if err != nil {
		return fmt.Errorf("while fuzzing (%q -> %q): %w", rc.inputPathOrEmpty(), rc.OutLitmus(), err)
	}
	return f.lift(ctx, rc, fuzzed)
}

func (f *FuzzRunner) check() error {
	if f.Fuzzer == nil {
		return ErrNoFuzzer
	}
	if f.Lifter == nil {
		return ErrNoLifter
	}
	if f.StatDumper == nil {
		return ErrNoStatDumper
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

func (f *FuzzRunner) fuzzJob(rc RunContext) (fuzzer2.Job, error) {
	in, err := rc.inputPath()
	if err != nil {
		return fuzzer2.Job{}, err
	}

	return fuzzer2.Job{
		Seed:      rc.Seed,
		In:        filepath.ToSlash(in),
		OutLitmus: rc.OutLitmus(),
		Config:    f.Config,
	}, nil
}

func (f *FuzzRunner) lift(ctx context.Context, rc RunContext, lpath string) error {
	if err := yos.MakeDir(rc.LiftOutDir()); err != nil {
		return fmt.Errorf("making dir for lift output: %w", err)
	}
	in, err := backend.InputFromFile(ctx, lpath, f.StatDumper)
	if err != nil {
		return fmt.Errorf("reading in input file %q: %w", lpath, err)
	}
	// TODO(@MattWindsor91): do something with the recipe
	_, err = f.Lifter.Lift(ctx, f.liftJob(in, rc), f.Runner)
	return err
}

func (f *FuzzRunner) liftJob(in backend.LiftInput, rc RunContext) backend.LiftJob {
	return backend.LiftJob{
		Arch: f.Arch,
		In:   in,
		Out: backend.LiftOutput{
			Dir:    rc.LiftOutDir(),
			Target: backend.ToDefault,
		},
	}
}
