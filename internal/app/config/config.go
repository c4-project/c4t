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

	"github.com/c4-project/c4t/internal/config/pretty"

	"github.com/mitchellh/go-wordwrap"

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

With no arguments, it produces an initial config file for the current system and dumps it to stdout.

The following flags print information about the current c4t config:

-` + FlagPrintGlobalPath + `: prints the path that c4t uses by default when looking for a config file.
You can use this to open the global config in a text editor, or save the config file produced by this program there.

-` + FlagPrintCompilers + `: prints a list of the currently-configured compilers.`

	// FlagPrintGlobalPath is the flag used for printing the global path.
	FlagPrintGlobalPath      = "print-global-path"
	flagPrintGlobalPathShort = "G"
	usagePrintPath           = "print path to global config file, rather than generating a new one"

	FlagPrintCompilers      = "print-compilers"
	flagPrintCompilersShort = "C"
	usagePrintCompilers     = "print information about configured compilers"
)

// App is the entry point for c4t-config.
func App(outw, errw io.Writer) *c.App {
	a := &c.App{
		Name:        Name,
		Usage:       usage,
		Description: wordwrap.WrapString(readme, 80),
		Flags:       flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetCommonAppSettings(a, outw, errw)
}

func flags() []c.Flag {
	ownFlags := []c.Flag{
		&c.BoolFlag{
			Name:    FlagPrintGlobalPath,
			Aliases: []string{flagPrintGlobalPathShort},
			Usage:   usagePrintPath,
		},
		&c.BoolFlag{
			Name:    FlagPrintCompilers,
			Aliases: []string{flagPrintCompilersShort},
			Usage:   usagePrintCompilers,
		},
	}
	return append(ownFlags, stdflag.TabulatorCliFlags()...)
}

func run(ctx *c.Context, outw io.Writer, errw io.Writer) error {
	if dc := dumpConfigFromCli(ctx); dc.isDumping() {
		return dump(ctx, outw, dc)
	}
	return probeAndDump(ctx.Context, outw, errw)
}

func dump(ctx *c.Context, outw io.Writer, dc *dumpConfig) error {
	var (
		// Loaded only if we need it for dumping.
		cfg *config.Config
		err error
	)
	if dc.needsConfig() {
		if cfg, err = stdflag.ConfigFromCli(ctx); err != nil {
			return err
		}
	}
	if dc.globalPath {
		if err = printGlobalPath(outw); err != nil {
			return err
		}
	}
	if dc.compilers {
		t := stdflag.TabulatorFromCli(ctx, outw)
		if err = pretty.TabulateCompilers(t, cfg); err != nil {
			return err
		}
	}
	return err
}

type dumpConfig struct {
	globalPath bool
	compilers  bool
}

func dumpConfigFromCli(ctx *c.Context) *dumpConfig {
	return &dumpConfig{
		globalPath: ctx.Bool(FlagPrintGlobalPath),
		compilers:  ctx.Bool(FlagPrintCompilers),
	}
}

func (d dumpConfig) isDumping() bool {
	// to be expanded if we add more config dumpers
	return d.compilers || d.globalPath
}

func (d dumpConfig) needsConfig() bool {
	// to be expanded if we add more config dumpers
	return d.compilers
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
