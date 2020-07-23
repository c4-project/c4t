// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package mach

import (
	"context"
	"io"
	"log"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/forward"

	"github.com/MattWindsor91/act-tester/internal/plan"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/compiler"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/runner"
)

// Config configures the machine-dependent stage.
type Config struct {
	// CDriver is the main driver for the compiler.
	CDriver compiler.SingleRunner
	// RDriver is the main driver for the runner.
	RDriver runner.ObsParser
	// Stdout is the Writer to which standard out from the machine-dependent stage should go.
	Stdout io.Writer
	// Logger is the logger to use for the machine-dependent stage.
	Logger *log.Logger
	// Observers is the set of observers to attach on the machine-dependent stage.
	Observers []builder.Observer
	// Json is, if present, a pointer to the JSON forwarding observer.
	Json *forward.Observer
	// User is the user configuration.
	User UserConfig
}

// UserConfig contains the part of the machine-stage configuration that can be set by the user,
// either directly or through invoker.
type UserConfig struct {
	// OutDir is the path to the output directory.
	OutDir string
	// SkipCompiler tells the machine-runner to skip compilation.
	SkipCompiler bool
	// SkipRunner tells the machine-runner to skip running.
	SkipRunner bool
	// Quantities contains various tunable quantities for the machine-dependent stage.
	Quantities QuantitySet
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
	if c.User.SkipCompiler {
		return nil
	}

	return &compiler.Config{
		Driver:     c.CDriver,
		Logger:     c.Logger,
		Paths:      compiler.NewPathset(c.User.OutDir),
		Observers:  c.Observers,
		Quantities: c.User.Quantities.Compiler,
	}
}

func (c *Config) makeRunnerConfig() *runner.Config {
	if c.User.SkipRunner {
		return nil
	}

	return &runner.Config{
		Logger:     c.Logger,
		Parser:     c.RDriver,
		Paths:      runner.NewPathset(c.User.OutDir),
		Observers:  c.Observers,
		Quantities: c.User.Quantities.Runner,
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
