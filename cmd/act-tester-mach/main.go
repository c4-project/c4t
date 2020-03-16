// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/pkg/mach/forward"

	"github.com/MattWindsor91/act-tester/internal/pkg/resolve"

	"github.com/MattWindsor91/act-tester/internal/pkg/mach/compiler"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/mach/runner"

	"github.com/MattWindsor91/act-tester/internal/pkg/act"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

const defaultOutDir = "mach_results"

func main() {
	if err := run(os.Args, os.Stdout, os.Stderr); err != nil {
		// TODO(@MattWindsor91): make this work properly with JSON output.
		ux.LogTopError(err)
	}
}

func run(args []string, outw, errw io.Writer) error {
	var pfile string
	a := act.Runner{Stderr: errw}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	c := makeConfigFlags(fs)

	ux.ActRunnerFlags(fs, &a)
	ux.PlanFileFlag(fs, &pfile)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	ccfg := c.makeCompilerConfig(errw)
	rcfg := c.makeRunnerConfig(errw, a)
	return runOnConfigs(context.Background(), ccfg, rcfg, pfile, outw)
}

type Config struct {
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

func makeConfigFlags(fs *flag.FlagSet) *Config {
	var c Config
	fs.BoolVar(&c.SkipCompiler, "c", false, "if given, skip the compiler")
	fs.BoolVar(&c.SkipRunner, "r", false, "if given, skip the runner")
	fs.IntVar(&c.Timeout, "t", 1, "a timeout, in `minutes`, to apply to each run")
	fs.IntVar(&c.NWorkers, "j", 1, "number of `workers` to run in parallel")
	fs.BoolVar(&c.JsonStatus, "J", false, "emit progress reports in JSON form on stderr")
	ux.OutDirFlag(fs, &c.OutDir, defaultOutDir)
	return &c
}

func (c *Config) makeCompilerConfig(errw io.Writer) *compiler.Config {
	if c.SkipCompiler {
		return nil
	}

	l, obs := c.makeLoggerAndObserver(errw, "compiler: ")
	return &compiler.Config{
		Driver:   &resolve.CResolve,
		Logger:   l,
		Paths:    compiler.NewPathset(c.OutDir),
		Observer: obs,
	}
}

func (c *Config) makeRunnerConfig(errw io.Writer, act act.Runner) *runner.Config {
	if c.SkipRunner {
		return nil
	}

	l, obs := c.makeLoggerAndObserver(errw, "runner: ")
	return &runner.Config{
		Logger:   l,
		Parser:   &act,
		Paths:    runner.NewPathset(c.OutDir),
		Observer: obs,
		Timeout:  c.Timeout,
		NWorkers: c.NWorkers,
	}
}

func (c *Config) makeLoggerAndObserver(errw io.Writer, prefix string) (*log.Logger, builder.Observer) {
	if c.JsonStatus {
		return nil, &forward.Observer{Encoder: json.NewEncoder(errw)}
	}
	l := log.New(errw, prefix, log.LstdFlags)
	return l, ux.NewPbObserver(l)
}

func runOnConfigs(ctx context.Context, cc *compiler.Config, rc *runner.Config, pfile string, outw io.Writer) error {
	p, perr := ux.LoadPlan(pfile)
	if perr != nil {
		return perr
	}
	cp, cerr := runCompiler(ctx, cc, p)
	if cerr != nil {
		return fmt.Errorf("while running compiler: %w", cerr)
	}
	rp, rerr := runRunner(ctx, rc, cp)
	if rerr != nil {
		return fmt.Errorf("while running runner: %w", rerr)
	}
	return rp.Dump(outw)
}

// runCompiler runs the batch compiler on plan p using config c, if available.
// If c is nil, runCompiler returns p unmodified.
func runCompiler(ctx context.Context, c *compiler.Config, p *plan.Plan) (*plan.Plan, error) {
	if c == nil {
		return p, nil
	}
	return c.Run(ctx, p)
}

// runRunner runs the batch runner on plan p using config c, if available.
// If c is nil, runRunner returns p unmodified.
func runRunner(ctx context.Context, c *runner.Config, p *plan.Plan) (*plan.Plan, error) {
	if c == nil {
		return p, nil
	}

	run, rerr := runner.New(c, p)
	if rerr != nil {
		return nil, rerr
	}
	out, oerr := run.Run(ctx)
	if oerr != nil {
		return nil, oerr
	}
	return out, nil
}
