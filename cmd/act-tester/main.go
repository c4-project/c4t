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
	"github.com/MattWindsor91/act-tester/internal/pkg/planner"

	"github.com/MattWindsor91/act-tester/internal/pkg/config"

	"github.com/MattWindsor91/act-tester/internal/pkg/director"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

func main() {
	err := run(os.Args, os.Stderr)
	ux.LogTopError(err)
}

func run(args []string, errw io.Writer) error {
	a := act.Runner{Stderr: errw}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	ux.ActRunnerFlags(fs, &a)
	cfile := ux.ConfFileFlag(fs)
	qs := setupQuantityOverrides(fs)

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	c, err := config.Load(*cfile)
	if err != nil {
		return err
	}
	c.Quantities.Override(*qs)

	e := makeEnv(&a, c)

	l := log.New(errw, "", 0)
	dc, err := director.ConfigFromGlobal(c, l, e)
	if err != nil {
		return nil
	}

	d, err := director.New(dc, fs.Args())
	if err != nil {
		return err
	}
	return d.Direct(context.Background())
}

func makeEnv(a *act.Runner, c *config.Config) director.Env {
	return director.Env{
		Fuzzer: a,
		Planner: planner.Source{
			BProbe: a,
			CProbe: c,
			SProbe: a,
		},
	}
}

func setupQuantityOverrides(fs *flag.FlagSet) *config.QuantitySet {
	var q config.QuantitySet
	// TODO(@MattWindsor91): disambiguate the corpus size argument
	ux.CorpusSizeFlag(fs, &q.Fuzz.CorpusSize)
	ux.SubjectCycleFlag(fs, &q.Fuzz.SubjectCycles)
	return &q
}
