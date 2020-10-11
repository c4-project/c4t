// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/1set/gut/ystring"
	"github.com/MattWindsor91/act-tester/internal/model/job"
	fuzzer2 "github.com/MattWindsor91/act-tester/internal/model/service/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/stage/fuzzer"
)

// ErrNoInput occurs when we try to run a mutating fuzzer for coverage and haven't any input to feed it.
var ErrNoInput = errors.New("this runner needs input testcases, but got none")

// FuzzRunner is a coverage runner that uses the act fuzzer.
type FuzzRunner struct {
	// Fuzzer is the fuzzer this fuzz runner uses.
	Fuzzer fuzzer.SingleFuzzer
	// Config is the configuration to pass to the fuzz runner.
	Config *fuzzer2.Configuration
}

func (f *FuzzRunner) Run(ctx context.Context, rc RunnerContext) error {
	j, err := f.job(rc)
	if err != nil {
		return fmt.Errorf("can't make fuzzer job: %w", err)
	}
	return f.Fuzzer.Fuzz(ctx, j)
}

func (f *FuzzRunner) job(rc RunnerContext) (job.Fuzzer, error) {
	input := rc.inputPath()
	if ystring.IsBlank(input) {
		return job.Fuzzer{}, ErrNoInput
	}

	return job.Fuzzer{
		Seed:      rc.Seed,
		In:        rc.inputPath(),
		OutLitmus: filepath.Join(rc.BucketDir, fmt.Sprintf("%d.litmus", rc.NumInBucket)),
		Config:    f.Config,
	}, nil
}
