// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package config implements the c4f-config app.
package config

import (
	"context"
	"fmt"
	"io"

	"github.com/c4-project/c4t/internal/helper/srvrun"

	"github.com/c4-project/c4t/internal/config"
	"github.com/c4-project/c4t/internal/machine"
	"github.com/c4-project/c4t/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	// Name is the name of the config binary.
	Name  = "c4t-config"
	usage = "manipulates config"

	readme = `
This program provides utilities for manipulating the c4t config file.

With no arguments, it produces an initial c4 config file for the current system
and dumps it to stdout.

With the -` + FlagPrintGlobalPath + ` argument, it prints the path that c4t uses by
default when looking for a config file.  You can use this feature to open the
global config in a text editor, or save the config file produced by this program
there.
`

	// FlagPrintGlobalPath is the flag used for printing the global path.
	FlagPrintGlobalPath      = "print-global-path"
	flagPrintGlobalPathShort = "G"
	usagePrintPath           = "print path to global config file, rather than generating a new one"
)

// App is the entry point for c4t-config.
func App(outw, errw io.Writer) *c.App {
	a := &c.App{
		Name:        Name,
		Usage:       usage,
		Description: readme,
		Flags:       flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetCommonAppSettings(a, outw, errw)
}

func flags() []c.Flag {
	return []c.Flag{
		&c.BoolFlag{
			Name:    FlagPrintGlobalPath,
			Aliases: []string{flagPrintGlobalPathShort},
			Usage:   usagePrintPath,
		},
	}
}

func run(ctx *c.Context, outw io.Writer, errw io.Writer) error {
	if ctx.Bool(FlagPrintGlobalPath) {
		return printGlobalPath(outw)
	}
	return probeAndDump(ctx.Context, outw, errw)
}

func printGlobalPath(outw io.Writer) error {
	path, err := config.GlobalFile()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(outw, path)
	return err
}

func probeAndDump(ctx context.Context, outw io.Writer, errw io.Writer) error {
	cfg := config.Config{}
	if err := cfg.Probe(ctx, srvrun.NewExecRunner(srvrun.StderrTo(errw)), machine.LocalProber{}); err != nil {
		return err
	}
	return cfg.Dump(outw)
}
