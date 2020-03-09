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

	"github.com/MattWindsor91/act-tester/internal/pkg/config"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"github.com/MattWindsor91/act-tester/internal/pkg/act"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"

	"github.com/MattWindsor91/act-tester/internal/pkg/planner"
)

const usageMach = "ID of machine to use for this test plan"

func main() {
	err := run(os.Args, os.Stdout, os.Stderr)
	ux.LogTopError(err)
}

func run(args []string, outw, errw io.Writer) error {
	a := act.Runner{Stderr: errw}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	pmach := fs.String(ux.FlagMachine, "", usageMach)
	ux.ActRunnerFlags(fs, &a)

	cfile := ux.ConfFileFlag(fs)

	var cs int
	ux.CorpusSizeFlag(fs, &cs)

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	plan, err := makePlanner(*cfile, errw, a, fs.Args(), *pmach)
	if err != nil {
		return err
	}

	p, err := plan.Plan(context.Background())
	if err != nil {
		return err
	}

	return p.Dump(outw)
}

func makePlanner(cfile string, errw io.Writer, a act.Runner, inFiles []string, midstr string) (*planner.Planner, error) {
	c, err := config.Load(cfile)
	if err != nil {
		return nil, err
	}
	mid, err := model.TryIDFromString(midstr)
	if err != nil {
		return nil, err
	}

	l := log.New(errw, "", 0)
	plan := planner.Planner{
		Source: planner.Source{
			BProbe: &a,
			CProbe: c,
			SProbe: &a,
		},
		Logger:    l,
		Observer:  ux.NewPbObserver(l),
		InFiles:   inFiles,
		MachineID: mid,
	}
	return &plan, nil
}
