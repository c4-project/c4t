// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/tools/gccnt"

	"github.com/MattWindsor91/act-tester/internal/pkg/ux"

	// This name is because every single time I try to use v2 named as 'cli', my IDE decides to replace it with v1.
	// Yes, I know, I shouldn't work around IDE issues by obfuscating my code, but I'm at my wit's end.
	c "github.com/urfave/cli/v2"
)

func main() {
	app := c.App{
		Name:                   "act-gccnt",
		Usage:                  "wraps gcc with various optional failure modes",
		Flags:                  flags(),
		Action:                 run,
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
	}
	ux.LogTopError(app.Run(os.Args))
}

const (
	flagOutput       = "o"
	flagBin          = "nt-bin"
	flagDryRun       = "nt-dryrun"
	flagDivergeOnOpt = "nt-diverge-opt"
	flagErrorOnOpt   = "nt-error-opt"
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
			Usage:   "O-levels (minus the '-O') on which gccn't should exit with an error",
		},
		&c.StringSliceFlag{
			Name:    flagDivergeOnOpt,
			Aliases: nil,
			Usage:   "O-levels (minus the '-O') on which gccn't should diverge",
		},
	}
	return append(fs, oflags()...)
}

var oflagNames = []string{"0", "1", "2", "3", "fast", "s", "g"}

func oflags() []c.Flag {
	flags := make([]c.Flag, len(oflagNames))
	for i, o := range oflagNames {
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
		Bin:         "gcc",
		In:          ctx.Args().Slice(),
		Out:         ctx.Path(flagOutput),
		OptLevel:    olevel,
		DivergeOpts: ctx.StringSlice(flagDivergeOnOpt),
		ErrorOpts:   ctx.StringSlice(flagErrorOnOpt),
	}

	if ctx.Bool(flagDryRun) {
		return g.DryRun(ctx.Context, os.Stderr)
	}
	return g.Run(ctx.Context, os.Stdout, os.Stderr)
}

func geto(ctx *c.Context) (string, error) {
	set := false
	o := "0"

	for _, possible := range oflagNames {
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
