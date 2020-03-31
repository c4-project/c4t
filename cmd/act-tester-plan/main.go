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

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/config"

	"github.com/MattWindsor91/act-tester/internal/act"
	"github.com/MattWindsor91/act-tester/internal/view"

	"github.com/MattWindsor91/act-tester/internal/controller/planner"
)

const usageMach = "ID of machine to use for this test plan"

func main() {
	err := run(os.Args, os.Stdout, os.Stderr)
	view.LogTopError(err)
}

func run(args []string, outw, errw io.Writer) error {
	a := act.Runner{Stderr: errw}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	pmach := fs.String(view.FlagMachine, "", usageMach)
	view.ActRunnerFlags(fs, &a)

	cfile := view.ConfFileFlag(fs)

	var cs int
	view.CorpusSizeFlag(fs, &cs)

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	plan, err := makePlanner(*cfile, errw, a, *pmach, cs)
	if err != nil {
		return err
	}

	p, err := plan.Plan(context.Background(), fs.Args())
	if err != nil {
		return err
	}

	return p.Dump(outw)
}

func makePlanner(cfile string, errw io.Writer, a act.Runner, midstr string, cs int) (*planner.Planner, error) {
	c, err := config.Load(cfile)
	if err != nil {
		return nil, err
	}
	mid, err := id.TryFromString(midstr)
	if err != nil {
		return nil, err
	}

	l := log.New(errw, "", 0)
	plan := planner.Planner{
		CorpusSize: cs,
		Source: planner.Source{
			BProbe: c,
			CProbe: c,
			SProbe: &a,
		},
		Logger:    l,
		Observers: view.Observers(l),
		MachineID: mid,
	}
	return &plan, nil
}
