// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package coverage contains the app definition for act-tester-coverage.
package coverage

import (
	"fmt"
	"io"
	"log"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"

	"github.com/MattWindsor91/act-tester/internal/ux/singleobs"

	"github.com/MattWindsor91/act-tester/internal/coverage"

	"github.com/MattWindsor91/act-tester/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	flagConfigFile      = "config"
	flagConfigFileShort = "c"
	usageConfigFile     = "Path to config file for coverage (not the tester config file!)"

	defaultConfigFile = "coverage.toml"
	defaultOutDir     = "coverage"
)

// App creates the act-tester-coverage app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "act-tester-coverage",
		Usage: "makes a coverage testbed using a plan",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, errw)
		},
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	ownFlags := []c.Flag{
		stdflag.VerboseFlag(),
		&c.StringFlag{
			Name:      flagConfigFile,
			Aliases:   []string{flagConfigFileShort},
			Usage:     usageConfigFile,
			TakesFile: true,
			Value:     defaultConfigFile,
		},
		stdflag.OutDirCliFlag(defaultOutDir),
	}
	return append(ownFlags, stdflag.ActRunnerCliFlags()...)
}

func run(ctx *c.Context, errw io.Writer) error {
	l := log.New(errw, "", log.LstdFlags)
	ccfg, err := coverage.LoadConfigFromFile(ctx.String(flagConfigFile))
	if err != nil {
		return fmt.Errorf("opening coverage config file: %w", err)
	}
	a := stdflag.ActRunnerFromCli(ctx, errw)
	cm, err := ccfg.MakeMaker(
		coverage.SendStderrTo(errw),
		coverage.UseFuzzer(a),
		coverage.UseStatDumper(a),
		coverage.UseLifter(&backend.BResolve),
		coverage.ObserveWith(singleobs.Coverage(l, stdflag.Verbose(ctx))...),
	)
	if err != nil {
		return fmt.Errorf("setting up the coverage maker: %w", err)
	}
	return cm.Run(ctx.Context)
}
