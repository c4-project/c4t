// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"context"
	"flag"
	"io"
	"os"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler"

	"github.com/MattWindsor91/act-tester/internal/controller/mach"

	"github.com/MattWindsor91/act-tester/internal/act"
	"github.com/MattWindsor91/act-tester/internal/view"
)

const defaultOutDir = "mach_results"

func main() {
	if err := run(os.Args, os.Stdout, os.Stderr); err != nil {
		// TODO(@MattWindsor91): make this work properly with JSON output.
		view.LogTopError(err)
	}
}

func run(args []string, outw, errw io.Writer) error {
	var pfile string
	a := act.Runner{Stderr: errw}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	c := makeConfigFlags(fs)
	c.Stdout = outw
	c.Stderr = errw
	c.RDriver = &backend.BResolve
	c.CDriver = &compiler.CResolve

	view.ActRunnerFlags(fs, &a)
	view.PlanFileFlag(fs, &pfile)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	return view.RunOnPlanFile(context.Background(), c, pfile, outw)
}

func makeConfigFlags(fs *flag.FlagSet) *mach.Config {
	var c mach.Config
	fs.BoolVar(&c.SkipCompiler, "c", false, "if given, skip the compiler")
	fs.BoolVar(&c.SkipRunner, "r", false, "if given, skip the runner")
	fs.IntVar(&c.Timeout, "t", 1, "a timeout, in `minutes`, to apply to each run")
	fs.IntVar(&c.NWorkers, "j", 1, "number of `workers` to run in parallel")
	fs.BoolVar(&c.JsonStatus, "J", false, "emit progress reports in JSON form on stderr")
	view.OutDirFlag(fs, &c.OutDir, defaultOutDir)
	return &c
}
