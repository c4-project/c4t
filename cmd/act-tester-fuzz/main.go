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

const (
	// defaultOutDir is the default directory used for the results of the lifter.
	defaultOutDir = "fuzz_results"

	usageSubjectCycles = "number of `cycles` to run for each subject in the corpus"
)

func main() {
	err := run(os.Args, os.Stderr)
	ux.LogTopError(err)
}

func run(args []string, errw io.Writer) error {
	a := act.Runner{Stderr: errw}
	l := log.New(errw, "", 0)
	cfg := fuzzer.Config{
		Driver:   &a,
		Observer: ux.NewPbObserver(l),
	}

	var dir, pf string
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	ux.ActRunnerFlags(fs, &a)
	ux.CorpusSizeFlag(fs, &cfg.CorpusSize)
	ux.OutDirFlag(fs, &dir, defaultOutDir)
	ux.PlanFileFlag(fs, &pf)
	fs.IntVar(&cfg.SubjectCycles, "k", fuzzer.DefaultSubjectCycles, usageSubjectCycles)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	cfg.Paths = fuzzer.NewPathset(dir)
	return ux.RunOnPlanFile(context.Background(), &cfg, pf)
}
