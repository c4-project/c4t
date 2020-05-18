// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package litmus contains the app definition for act-litmus.
package litmus

import (
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/tool/litmus"
	"github.com/MattWindsor91/act-tester/internal/view/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	// These flags deliberately line up with those of Litmus.
	flagC11     = "c11"
	flagCArch   = "carch"
	flagOutDir  = "o" // different from stdflag.FlagOutDir
	flagVerbose = "v"

	usageC11     = "for Litmus compatibility; ignored"
	usageCArch   = "C architecture to pass through to litmus"
	usageOutDir  = "output directory for harness"
	usageVerbose = "be more verbose"
)

// App creates the act-litmus app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "act-litmus",
		Usage: "wraps litmus with various issue workarounds",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, errw)
		},
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	fs := []c.Flag{
		&c.StringFlag{
			Name:  flagC11,
			Usage: usageC11,
		},
		&c.BoolFlag{
			Name:  flagVerbose,
			Usage: usageVerbose,
			Value: false,
		},
		&c.StringFlag{
			Name:  flagCArch,
			Usage: usageCArch,
		},
		&c.PathFlag{
			Name:  flagOutDir,
			Usage: usageOutDir,
		},
	}
	return append(fs, stdflag.ActRunnerCliFlags()...)
}

func run(ctx *c.Context, errw io.Writer) error {
	lit, err := makeLitmus(ctx, errw)
	if err != nil {
		return err
	}
	return lit.Run(ctx.Context)
}

// makeLitmus makes a litmus runner using the arguments in ctx and the standard error writer errw.
func makeLitmus(ctx *c.Context, errw io.Writer) (*litmus.Litmus, error) {
	anons := ctx.Args().Slice()
	if len(anons) != 1 {
		return nil, fmt.Errorf("expected precisely one anonymous argument; got %v", anons)
	}

	a := stdflag.ActRunnerFromCli(ctx, errw)
	cfg := litmus.Litmus{
		Stat:    a,
		CArch:   ctx.String(flagCArch),
		Verbose: ctx.Bool(flagVerbose),
		Pathset: litmus.Pathset{FileIn: anons[0], DirOut: ctx.Path(flagOutDir)},
		Fixset:  litmus.Fixset{},
		Err:     errw,
	}

	cfg.Pathset.FileIn = anons[0]
	return &cfg, nil
}
