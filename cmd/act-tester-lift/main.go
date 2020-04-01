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

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"

	"github.com/MattWindsor91/act-tester/internal/controller/lifter"
	"github.com/MattWindsor91/act-tester/internal/view"
)

// defaultOutDir is the default directory used for the results of the lifter.
const defaultOutDir = "lift_results"

func main() {
	err := run(os.Args, os.Stdout, os.Stderr)
	view.LogTopError(err)
}

func run(args []string, outw, errw io.Writer) error {
	var pf string
	l := log.New(errw, "", 0)

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	var od string
	view.OutDirFlag(fs, &od, defaultOutDir)
	view.PlanFileFlag(fs, &pf)

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	cfg := lifter.Config{
		Maker:     &backend.BResolve,
		Logger:    l,
		Observers: view.BuilderObservers(l),
		Paths:     lifter.NewPathset(od),
		Stderr:    errw,
	}

	return view.RunOnPlanFile(context.Background(), &cfg, pf, outw)
}
