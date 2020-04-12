// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/MattWindsor91/act-tester/internal/config"
	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/controller/rmach"

	c "github.com/urfave/cli/v2"

	"github.com/MattWindsor91/act-tester/internal/view/singleobs"

	"github.com/MattWindsor91/act-tester/internal/view"
)

const (
	defaultOutDir = "rmach_results"
)

func main() {
	app := c.App{
		Name:                   "act-tester-rmach",
		Usage:                  "runs the machine-dependent phase of an ACT test, potentially remotely",
		Flags:                  flags(),
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
		Action: func(ctx *c.Context) error {
			return run(ctx, os.Stdout, os.Stderr)
		},
	}
	view.LogTopError(app.Run(os.Args))
}

func flags() []c.Flag {
	return []c.Flag{
		// TODO(@MattWindsor91): forward some flags from rmach to mach
		view.ConfFileCliFlag(),
		view.OutDirCliFlag(defaultOutDir),
		view.PlanFileCliFlag(),
	}
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	cfg, err := view.ConfFileFromCli(ctx)
	if err != nil {
		return err
	}

	errw = iohelp.EnsureWriter(errw)
	rcfg := makeConfig(ctx, cfg, errw)
	pfile := view.PlanFileFromCli(ctx)
	return view.RunOnPlanFile(context.Background(), rcfg, pfile, outw)
}

func makeConfig(ctx *c.Context, cfg *config.Config, errw io.Writer) *rmach.Config {
	l := log.New(errw, "[rmach] ", log.LstdFlags)
	obs := rmach.NewObserverSet(singleobs.RMach(l)...)
	return &rmach.Config{
		DirLocal:  view.OutDirFromCli(ctx),
		Observers: obs,
		SSH:       cfg.SSH,
	}
}
