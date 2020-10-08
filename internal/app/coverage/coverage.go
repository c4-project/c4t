// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package coverage contains the app definition for act-tester-coverage.
package coverage

import (
	"io"
	"os"

	"github.com/MattWindsor91/act-tester/internal/coverage"
	"github.com/MattWindsor91/act-tester/internal/ux"

	"github.com/MattWindsor91/act-tester/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	flagConfigFile      = "config"
	flagConfigFileShort = "c"
	usageConfigFile     = "Path to config file for coverage (not the tester config file!)"

	defaultOutDir = "coverage"
)

// App creates the act-tester-coverage app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "act-tester-coverage",
		Usage: "makes a coverage testbed using a plan",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, os.Stdout, os.Stderr)
		},
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	return []c.Flag{
		//stdflag.VerboseFlag(),
		&c.StringFlag{
			Name:      flagConfigFile,
			Aliases:   []string{flagConfigFileShort},
			Usage:     usageConfigFile,
			TakesFile: true,
		},
		stdflag.OutDirCliFlag(defaultOutDir),
	}
}

func run(ctx *c.Context, outw, _errw io.Writer) error {
	ccfg, err := coverage.LoadConfigFromFile(ctx.String(flagConfigFile))
	if err != nil {
		return err
	}

	cm, err := coverage.NewMaker(ccfg.Profiles, coverage.OptionsFromConfig(ccfg))
	if err != nil {
		return err
	}
	return ux.RunOnCliPlan(ctx, cm, outw)
}
