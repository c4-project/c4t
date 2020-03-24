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

	"github.com/MattWindsor91/act-tester/internal/pkg/ux/dash"

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

	return runWithArgs(cfile, qs, a, fs.Args())
}

func runWithArgs(cfile *string, qs *config.QuantitySet, a act.Runner, files []string) error {
	c, err := config.Load(*cfile)
	if err != nil {
		return err
	}
	dc, err := makeDirectorConfig(c, qs, a)
	if err != nil {
		return err
	}

	d, err := director.New(dc, files)
	if err != nil {
		return err
	}
	return d.Direct(context.Background())
}

func makeDirectorConfig(c *config.Config, qs *config.QuantitySet, a act.Runner) (*director.Config, error) {
	c.Quantities.Override(*qs)

	// TODO(@MattWindsor91)
	e := makeEnv(&a, c)
	mids, err := c.MachineIDs()
	if err != nil {
		return nil, err
	}
	o, err := dash.New(mids)
	if err != nil {
		return nil, err
	}
	l := log.New(o, "", 0)
	dc, err := director.ConfigFromGlobal(c, l, e, o)
	if err != nil {
		return nil, err
	}
	return dc, nil
}

func makeEnv(a *act.Runner, c *config.Config) director.Env {
	return director.Env{
		Fuzzer: a,
		Lifter: &backend.BResolve,
		Planner: planner.Source{
			BProbe: c,
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
