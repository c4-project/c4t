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

const usageConfFile = "The `file` from which to load the tester configuration."

func run(args []string, errw io.Writer) error {
	a := act.Runner{Stderr: errw}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	ux.ActRunnerFlags(fs, &a)
	cfile := fs.String("C", "", usageConfFile)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	c, err := config.Load(*cfile)
	if err != nil {
		return err
	}

	e := director.Env{Planner: planner.Source{
		BProbe: &a,
		CProbe: c,
		SProbe: &a}}

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
