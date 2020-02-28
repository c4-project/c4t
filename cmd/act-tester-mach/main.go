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
	ux.ActRunnerFlags(fs, &act)
	ux.OutDirFlag(fs, &dir, defaultOutDir)
	ux.PlanFileFlag(fs, &pfile)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	ccfg := compiler.Config{
		Driver: &act,
		Logger: log.New(errw, "compiler: ", 0),
		Paths:  compiler.NewPathset(dir),
	}
	rcfg := runner.Config{
		Logger: log.New(errw, "runner: ", 0),
		Parser: &act,
		Paths:  runner.NewPathset(dir),
	}
	return runOnConfigs(context.Background(), &ccfg, &rcfg, pfile, outw)
}

func runOnConfigs(ctx context.Context, cc *compiler.Config, rc *runner.Config, pfile string, outw io.Writer) error {
	p, perr := ux.LoadPlan(pfile)
	if perr != nil {
		return perr
	}
	cp, cerr := cc.Run(ctx, p)
	if cerr != nil {
		return fmt.Errorf("while running compiler: %w", cerr)
	}
	rres, rerr := runRunner(ctx, rc, cp)
	if rerr != nil {
		return fmt.Errorf("while running runner: %w", rerr)
	}
	// TODO(@MattWindsor91): filter
	return json.NewEncoder(outw).Encode(rres)
}

func runRunner(ctx context.Context, c *runner.Config, p *plan.Plan) (*runner.Result, error) {
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
