// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package mach

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/MattWindsor91/act-tester/internal/controller/mach/forward"

	"github.com/MattWindsor91/act-tester/internal/model/plan"

	"github.com/MattWindsor91/act-tester/internal/controller/mach/compiler"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/runner"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
)

// Config configures the machine-dependent stage.
type Config struct {
	// CDriver is the main driver for the compiler.
	CDriver compiler.SingleRunner
	// RDriver is the main driver for the runner.
	RDriver runner.ObsParser
	// Stdout is the Writer to which standard out from the machine-dependent stage should go.
	Stdout io.Writer
	// OutDir is the path to the output directory.
	OutDir string
	// SkipCompiler tells the machine-runner to skip compilation.
	SkipCompiler bool
	// SkipRunner tells the machine-runner to skip running.
	SkipRunner bool
	// CTimeout is a timeout to apply to each compilation.
	CTimeout time.Duration
	// RTimeout is a timeout to apply to each run.
	RTimeout time.Duration
	// NWorkers is the number of running workers to run in parallel.
	NWorkers int
	// Logger is the logger to use for the machine-dependent stage.
	Logger *log.Logger
	// Observers is the set of observers to attach on the machine-dependent stage.
	Observers []builder.Observer
	// Json is, if present, a pointer to the JSON forwarding observer.
	Json *forward.Observer
}

func (c *Config) Check() error {
	if c.CDriver == nil {
		return compiler.ErrDriverNil
	}
	if c.RDriver == nil {
		return runner.ErrParserNil
	}
	return nil
}

func (c *Config) makeCompilerConfig() *compiler.Config {
	if c.SkipCompiler {
		return nil
	}

	return &compiler.Config{
		Driver:    c.CDriver,
		Logger:    c.Logger,
		Paths:     compiler.NewPathset(c.OutDir),
		Observers: c.Observers,
		Timeout:   c.CTimeout,
	}
}

func (c *Config) makeRunnerConfig() *runner.Config {
	if c.SkipRunner {
		return nil
	}

	return &runner.Config{
		Logger:    c.Logger,
		Parser:    c.RDriver,
		Paths:     runner.NewPathset(c.OutDir),
		Observers: c.Observers,
		Timeout:   c.RTimeout,
		NWorkers:  c.NWorkers,
	}
}

// Run creates a new machine-dependent phase runner from this config, then runs it on p using ctx.
func (c *Config) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	m, err := New(c, p)
	if err != nil {
		return nil, err
	}
	return m.Run(ctx)
}
