// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/compiler"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/runner"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

const defaultOutDir = "mach_results"

func main() {
	if err := run(os.Args, os.Stdout, os.Stderr); err != nil {
		ux.LogTopError(err)
	}
}

func run(args []string, outw, errw io.Writer) error {
	var (
		dir   string
		pfile string
	)
	act := interop.ActRunner{Stderr: errw}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	skipc := fs.Bool("c", false, "If given, skip the compiler")
	skipr := fs.Bool("r", false, "If given, skip the runner")
	ux.ActRunnerFlags(fs, &act)
	ux.OutDirFlag(fs, &dir, defaultOutDir)
	ux.PlanFileFlag(fs, &pfile)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	ccfg := makeCompilerConfig(*skipc, errw, act, dir)
	rcfg := makeRunnerConfig(*skipr, errw, act, dir)
	return runOnConfigs(context.Background(), ccfg, rcfg, pfile, outw)
}

func makeCompilerConfig(skip bool, errw io.Writer, act interop.ActRunner, dir string) *compiler.Config {
	if skip {
		return nil
	}

	cl := log.New(errw, "compiler: ", 0)
	return &compiler.Config{
		Driver:   &act,
		Logger:   cl,
		Paths:    compiler.NewPathset(dir),
		Observer: ux.NewPbObserver(cl),
	}
}

func makeRunnerConfig(skip bool, errw io.Writer, act interop.ActRunner, dir string) *runner.Config {
	if skip {
		return nil
	}

	rl := log.New(errw, "runner: ", 0)
	return &runner.Config{
		Logger:   rl,
		Parser:   &act,
		Paths:    runner.NewPathset(dir),
		Observer: ux.NewPbObserver(rl),
	}
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
