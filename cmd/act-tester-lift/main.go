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

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/lifter"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

// defaultOutDir is the default directory used for the results of the lifter.
const defaultOutDir = "lift_results"

func main() {
	err := run(os.Args, os.Stderr)
	ux.LogTopError(err)
}

func run(args []string, errw io.Writer) error {
	var pf string
	act := interop.ActRunner{Stderr: errw}
	l := log.New(errw, "", 0)
	lift := lifter.Lifter{
		Maker:    &act,
		Observer: ux.NewPbObserver(l),
	}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	ux.ActRunnerFlags(fs, &act)
	ux.OutDirFlag(fs, &lift.OutDir, defaultOutDir)
	ux.PlanFileFlag(fs, &pf)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	return ux.RunOnPlanFile(context.Background(), &lift, pf)
}
