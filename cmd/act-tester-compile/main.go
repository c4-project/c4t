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

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/compiler"
	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

const (
	defaultOutDir = "compile_results"
)

func main() {
	if err := run(os.Args, os.Stderr); err != nil {
		ux.LogTopError(err)
	}
}

func run(args []string, errw io.Writer) error {
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

	cfg := compiler.Config{
		Driver:   &act,
		Logger:   log.New(errw, "", 0),
		Paths:    compiler.NewPathset(dir),
		Observer: &corpus.PbObserver{},
	}
	return ux.RunOnPlanFile(context.Background(), &cfg, pfile)
}
