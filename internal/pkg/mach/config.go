// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package mach

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/pkg/mach/forward"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"

	"github.com/MattWindsor91/act-tester/internal/pkg/mach/compiler"
	"github.com/MattWindsor91/act-tester/internal/pkg/mach/runner"
)

// Config configures the machine-dependent stage.
type Config struct {
	// CDriver is the main driver for the compiler.
	CDriver compiler.SingleRunner
	// RDriver is the main driver for the runner.
	RDriver runner.ObsParser
	// Stdout is the Writer to which standard out from the machine-dependent stage should go.
	Stdout io.Writer
	// Stderr is the Writer to which standard error from the machine-dependent stage should go.
	Stderr io.Writer
	// OutDir is the path to the output directory.
	OutDir string
	// SkipCompiler tells the machine-runner to skip compilation.
	SkipCompiler bool
	// SkipRunner tells the machine-runner to skip running.
	SkipRunner bool
	// Timeout is a timeout, in minutes, to apply to each run.
	Timeout int
	// NWorkers is the number of running workers to run in parallel.
	NWorkers int
	// JsonStatus switches the machine-runner from human-readable status output to JSON status output.
	JsonStatus bool
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

	l, obs := c.makeLoggerAndObserver("compiler: ")
	return &compiler.Config{
		Driver:   c.CDriver,
		Logger:   l,
		Paths:    compiler.NewPathset(c.OutDir),
		Observer: obs,
	}
}

func (c *Config) makeRunnerConfig() *runner.Config {
	if c.SkipRunner {
		return nil
	}

	l, obs := c.makeLoggerAndObserver("runner: ")
	return &runner.Config{
		Logger:   l,
		Parser:   c.RDriver,
		Paths:    runner.NewPathset(c.OutDir),
		Observer: obs,
		Timeout:  c.Timeout,
		NWorkers: c.NWorkers,
	}
}

func (c *Config) makeLoggerAndObserver(prefix string) (*log.Logger, builder.Observer) {
	errw := c.Stderr
	if errw == nil {
		errw = ioutil.Discard
	}

	if c.JsonStatus {
		return nil, &forward.Observer{Encoder: json.NewEncoder(errw)}
	}

	l := log.New(errw, prefix, log.LstdFlags)
	return l, ux.NewPbObserver(l)
}

// Run creates a new machine-dependent phase runner from this config, then runs it on p using ctx.
func (c *Config) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	m, err := New(c, p)
	if err != nil {
		return nil, err
	}
	return m.Run(ctx)
}
