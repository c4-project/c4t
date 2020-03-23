// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"context"
	"io"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/helpers/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/obs"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// ObsParser is the interface of things that can parse test outcomes.
type ObsParser interface {
	// ParseObs parses the observation in reader r into o according to the backend configuration in b.
	// The backend described by b must have been used to produce the testcase outputting r.
	ParseObs(ctx context.Context, b model.Backend, r io.Reader, o *obs.Obs) error
}

// MachConfig represents the configuration needed to run a Runner.
type Config struct {
	// Timeout is the timeout for each run, in minutes.
	// Non-positive values disable the timeout.
	Timeout int

	// NWorkers is the number of parallel run workers that should be spawned.
	// Anything less than or equal to 1 will sequentialise the run.
	NWorkers int

	// Logger is the logger that should be used for this Runner.
	// If nil, logging will be suppressed.
	Logger *log.Logger

	// Observer observes the runner's progress across a corpus.
	Observer builder.Observer

	// Parser handles the parsing of observations.
	Parser ObsParser

	// Paths contains the pathset used for this runner's outputs.
	Paths *Pathset
}

// Check checks various error conditions on the config.
func (c *Config) Check() error {
	if c.Parser == nil {
		return ErrParserNil
	}
	if c.Paths == nil {
		return iohelp.ErrPathsetNil
	}
	return nil
}

// Run constructs a new runner using this configuration, then runs it in ctx with p.
func (c *Config) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	run, rerr := New(c, p)
	if rerr != nil {
		return nil, rerr
	}
	out, oerr := run.Run(ctx)
	if oerr != nil {
		return nil, oerr
	}
	return out, nil
}
