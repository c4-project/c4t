// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"context"
	"log"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Config contains planner configuration that can persist across multiple plans on multiple machines.
type Config struct {
	// Source contains all of the various sources for a Planner's information.
	Source Source

	// Filter is the compiler filter to use to select compilers to test.
	Filter string

	// Logger is the logger used by the planner.
	Logger *log.Logger

	// Observers contains the set of observers used to get feedback on the planning action as it completes.
	Observers ObserverSet

	// CorpusSize is the requested size of the test corpus.
	// If zero, no corpus sampling is done, but the planner will still error if the final corpus size is 0.
	// If nonzero, the corpus will be sampled if larger than the size, and an error occurs if the final size is below
	// that requested.
	CorpusSize int
}

// Plan constructs a Planner using this config, then runs it using ctx on the file set fs and machine mach.
func (c *Config) Plan(ctx context.Context, mach plan.NamedMachine, fs []string, seed int64) (*plan.Plan, error) {
	p, err := New(c, mach, fs, seed)
	if err != nil {
		return nil, err
	}
	return p.Plan(ctx)
}

// Check performs in-flight checks on this config before its use.
func (c *Config) Check() error {
	// TODO(@MattWindsor91)
	return nil
}
