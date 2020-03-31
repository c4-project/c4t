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

	"github.com/MattWindsor91/act-tester/internal/act"
	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/view"
)

// defaultOutDir is the default directory used for the results of the lifter.
const defaultOutDir = "fuzz_results"

func main() {
	err := run(os.Args, os.Stdout, os.Stderr)
	view.LogTopError(err)
}

func run(args []string, outw, errw io.Writer) error {
	a := act.Runner{Stderr: errw}
	l := log.New(errw, "", 0)

	var dir, pf string
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	view.ActRunnerFlags(fs, &a)
	view.OutDirFlag(fs, &dir, defaultOutDir)
	view.PlanFileFlag(fs, &pf)

	qs := setupQuantityFlags(fs)

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	cfg := fuzzer.Config{
		Driver:     &a,
		Observers:  view.Observers(l),
		Logger:     l,
		Paths:      fuzzer.NewPathset(dir),
		Quantities: *qs,
	}
	return view.RunOnPlanFile(context.Background(), &cfg, pf, outw)
}

func setupQuantityFlags(fs *flag.FlagSet) *fuzzer.QuantitySet {
	var q fuzzer.QuantitySet
	view.CorpusSizeFlag(fs, &q.CorpusSize)
	view.SubjectCycleFlag(fs, &q.SubjectCycles)
	return &q
}
