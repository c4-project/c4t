// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package coverage contains the app definition for c4t-coverage.
package coverage

import (
	"fmt"
	"io"
	"log"

	"github.com/c4-project/c4t/internal/serviceimpl/backend"

	"github.com/c4-project/c4t/internal/ux/singleobs"

	"github.com/c4-project/c4t/internal/coverage"

	"github.com/c4-project/c4t/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	flagConfigFile      = "config"
	flagConfigFileShort = "c"
	usageConfigFile     = "Path to config file for coverage (not the tester config file!)"

	defaultConfigFile = "coverage.toml"
	defaultOutDir     = "coverage"
)

// App creates the c4t-coverage app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "c4t-coverage",
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
	return append(ownFlags, stdflag.C4fRunnerCliFlags()...)
}

func run(ctx *c.Context, errw io.Writer) error {
	l := log.New(errw, "", log.LstdFlags)
	ccfg, err := coverage.LoadConfigFromFile(ctx.String(flagConfigFile))
	if err != nil {
		return fmt.Errorf("opening coverage config file: %w", err)
	}
	a := stdflag.C4fRunnerFromCli(ctx, errw)
	cm, err := ccfg.MakeMaker(
		coverage.SendStderrTo(errw),
		coverage.UseFuzzer(a),
		coverage.UseStatDumper(a),
		coverage.UseBackendResolver(&backend.Resolve),
		coverage.ObserveWith(singleobs.Coverage(l, stdflag.Verbose(ctx))...),
	)
	if err != nil {
		return fmt.Errorf("setting up the coverage maker: %w", err)
	}
	return cm.Run(ctx.Context)
}
