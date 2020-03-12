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

	"github.com/MattWindsor91/act-tester/internal/pkg/act"
	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

// defaultOutDir is the default directory used for the results of the lifter.
const defaultOutDir = "fuzz_results"

func main() {
	err := run(os.Args, os.Stderr)
	ux.LogTopError(err)
}

func run(args []string, errw io.Writer) error {
	a := act.Runner{Stderr: errw}
	l := log.New(errw, "", 0)

	var dir, pf string
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	ux.ActRunnerFlags(fs, &a)
	ux.OutDirFlag(fs, &dir, defaultOutDir)
	ux.PlanFileFlag(fs, &pf)

	qs := setupQuantityFlags(fs)

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	cfg := fuzzer.Config{
		Driver:     &a,
		Observer:   ux.NewPbObserver(l),
		Logger:     l,
		Paths:      fuzzer.NewPathset(dir),
		Quantities: *qs,
	}
	return ux.RunOnPlanFile(context.Background(), &cfg, pf)
}

func setupQuantityFlags(fs *flag.FlagSet) *fuzzer.QuantitySet {
	var q fuzzer.QuantitySet
	ux.CorpusSizeFlag(fs, &q.CorpusSize)
	ux.SubjectCycleFlag(fs, &q.SubjectCycles)
	return &q
}
