// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/runner"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

const defaultOutDir = "run_results"

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

	cfg := runner.Config{
		Logger: log.New(errw, "", 0),
		Parser: &act,
		Paths:  runner.NewPathset(dir),
	}
	return makeAndRunRunner(&cfg, pfile, outw)
}

func makeAndRunRunner(c *runner.Config, pfile string, outw io.Writer) error {
	p, perr := ux.LoadPlan(pfile)
	if perr != nil {
		return perr
	}
	run, rerr := runner.New(c, p)
	if rerr != nil {
		return rerr
	}
	out, oerr := run.Run(context.Background())
	if oerr != nil {
		return oerr
	}
	return out.Dump(outw)
}
