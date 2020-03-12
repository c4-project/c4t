// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/act"

	"github.com/MattWindsor91/act-tester/internal/pkg/litmus"

	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

const (
	usageCArch   = "C architecture to pass through to litmus"
	usageOutDir  = "output directory for harness"
	usageVerbose = "be more verbose"
)

func main() {
	if err := run(os.Args, os.Stderr); err != nil {
		ux.LogTopError(err)
	}
}

// run contains the top-level logic of the command, passing through arguments args and stderr writer errw.
func run(args []string, errw io.Writer) error {
	cfg, err := parseArgs(args, errw)
	if err != nil {
		return err
	}

	return cfg.Run(context.Background())
}

// parseArgs parses args into a litmus config, using errw as stderr.
func parseArgs(args []string, errw io.Writer) (*litmus.Litmus, error) {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.SetOutput(errw)

	a := act.Runner{Stderr: errw}
	cfg := litmus.Litmus{Stat: &a, Err: errw}

	_ = fs.String("c11", "", "for Litmus compatibility; ignored")

	fs.BoolVar(&cfg.Verbose, "v", false, usageVerbose)
	fs.StringVar(&cfg.CArch, "carch", "", usageCArch)
	fs.StringVar(&cfg.Pathset.DirOut, "o", "", usageOutDir)
	ux.ActRunnerFlags(fs, &a)

	if err := fs.Parse(args[1:]); err != nil {
		return nil, err
	}

	anons := fs.Args()
	if len(anons) != 1 {
		return nil, fmt.Errorf("expected precisely one anonymous argument; got %v", anons)
	}
	cfg.Pathset.FileIn = anons[0]
	return &cfg, nil
}
