// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package gccnt contains the app definition for c4t-gccnt.
package gccnt

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/MattWindsor91/c4t/internal/ux/stdflag"

	"github.com/MattWindsor91/c4t/internal/serviceimpl/compiler/gcc"

	"github.com/MattWindsor91/c4t/internal/tool/gccnt"

	// This name is because every single time I try to use v2 named as 'cli', my IDE decides to replace it with v1.
	// Yes, I know, I shouldn't work around IDE issues by obfuscating my code, but I'm at my wit's end.
	c "github.com/urfave/cli/v2"
)

// App creates the c4t-gccnt app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:   "c4t-gccnt",
		Usage:  "wraps gcc with various optional failure modes",
		Flags:  flags(),
		Action: run,
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

const (
	flagOutput       = "o"
	flagBin          = "nt-bin"
	flagDryRun       = "nt-dryrun"
	flagDivergeOnOpt = "nt-diverge-opt"
	flagErrorOnOpt   = "nt-error-opt"
	flagPthread      = "pthread"
	flagStd          = "std"
	flagMarch        = "march"
	flagMcpu         = "mcpu"
)

func flags() []c.Flag {
	fs := []c.Flag{
		&c.PathFlag{
			Name:      flagOutput,
			Usage:     "output file",
			Required:  true,
			TakesFile: true,
			Value:     "a.out",
		},
		&c.StringFlag{
			Name:  flagBin,
			Usage: "the 'real' compiler `command` to run",
			Value: "gcc",
		},
		&c.BoolFlag{
			Name:  flagDryRun,
			Usage: "print the outcome of running gccn't instead of doing it",
		},
		&c.StringSliceFlag{
			Name:    flagErrorOnOpt,
			Aliases: nil,
			Usage:   "o-levels (minus the '-O') on which gccn't should exit with an error",
		},
		&c.StringSliceFlag{
			Name:    flagDivergeOnOpt,
			Aliases: nil,
			Usage:   "o-levels (minus the '-O') on which gccn't should diverge",
		},
		&c.StringFlag{
			Name:  flagStd,
			Usage: "standard to pass through to gcc",
		},
		&c.StringFlag{
			Name:  flagMarch,
			Usage: "architecture optimisation to pass through to gcc",
		},
		&c.StringFlag{
			Name:  flagMcpu,
			Usage: "cpu optimisation to pass through to gcc",
		},
		&c.BoolFlag{
			Name:  flagPthread,
			Usage: "passes through pthread to gcc",
		},
	}
	return append(fs, oflags()...)
}

func oflags() []c.Flag {
	flags := make([]c.Flag, len(gcc.OptLevelNames))
	for i, o := range gcc.OptLevelNames {
		flags[i] = &c.BoolFlag{
			Name:  "O" + o,
			Usage: fmt.Sprintf("optimisation level '%s'", o),
		}
	}
	return flags
}

func run(ctx *c.Context) error {
	olevel, err := geto(ctx)
	if err != nil {
		return err
	}

	g := gccnt.Gccnt{
		Bin:         ctx.String(flagBin),
		In:          ctx.Args().Slice(),
		Out:         ctx.Path(flagOutput),
		OptLevel:    olevel,
		DivergeOpts: ctx.StringSlice(flagDivergeOnOpt),
		ErrorOpts:   ctx.StringSlice(flagErrorOnOpt),
		March:       ctx.String(flagMarch),
		Mcpu:        ctx.String(flagMcpu),
		Pthread:     ctx.Bool(flagPthread),
		Std:         ctx.String(flagStd),
	}

	if ctx.Bool(flagDryRun) {
		return g.DryRun(ctx.Context, os.Stderr)
	}
	return g.Run(ctx.Context, os.Stdout, os.Stderr)
}

func geto(ctx *c.Context) (string, error) {
	set := false
	o := "0"

	for _, possible := range gcc.OptLevelNames {
		if ctx.Bool("O" + possible) {
			o = possible
			if set {
				return "", errors.New("multiple optimisation levels defined")
			}
			set = true
		}
	}

	return o, nil
}
