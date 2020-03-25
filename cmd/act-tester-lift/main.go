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

	"github.com/MattWindsor91/act-tester/internal/pkg/resolve/backend"

	"github.com/MattWindsor91/act-tester/internal/pkg/lifter"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

// defaultOutDir is the default directory used for the results of the lifter.
const defaultOutDir = "lift_results"

func main() {
	err := run(os.Args, os.Stdout, os.Stderr)
	ux.LogTopError(err)
}

func run(args []string, outw, errw io.Writer) error {
	var pf string
	l := log.New(errw, "", 0)

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	var od string
	ux.OutDirFlag(fs, &od, defaultOutDir)
	ux.PlanFileFlag(fs, &pf)

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	cfg := lifter.Config{
		Maker:    &backend.BResolve,
		Logger:   l,
		Observer: ux.NewPbObserver(l),
		Paths:    lifter.NewPathset(od),
		Stderr:   errw,
	}

	return ux.RunOnPlanFile(context.Background(), &cfg, pf, outw)
}
