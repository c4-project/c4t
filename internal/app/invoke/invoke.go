// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package invoker contains the app definition for c4t-invoke.
package invoke

import (
	"io"
	"log"

	"github.com/c4-project/c4t/internal/stage/invoker/runner"

	"github.com/c4-project/c4t/internal/helper/errhelp"

	"github.com/c4-project/c4t/internal/ux/stdflag"

	"github.com/c4-project/c4t/internal/config"
	"github.com/c4-project/c4t/internal/helper/iohelp"

	"github.com/c4-project/c4t/internal/stage/invoker"

	c "github.com/urfave/cli/v2"

	"github.com/c4-project/c4t/internal/ux/singleobs"

	"github.com/c4-project/c4t/internal/ux"
)

const (
	Name = "c4t-invoke"

	flagForce      = "force"
	flagForceShort = "f"
	usageForce     = "allow invoke on plans that have already been invoked"
)

// App creates the c4t-invoke app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  Name,
		Usage: "runs the machine-dependent phase of a C4 test, potentially remotely",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetPlanAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	ownFlags := []c.Flag{
		stdflag.VerboseFlag(),
		stdflag.ConfFileCliFlag(),
		&c.BoolFlag{
			Name:    flagForce,
			Aliases: []string{flagForceShort},
			Usage:   usageForce,
		},
	}
	return append(ownFlags, stdflag.MachCliFlags()...)
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	cfg, err := stdflag.ConfigFromCli(ctx)
	if err != nil {
		return err
	}

	inv, err := makeInvoker(ctx, cfg, iohelp.EnsureWriter(errw))
	if err != nil {
		return err
	}

	err = ux.RunOnCliPlan(ctx, inv, outw)
	cerr := inv.Close()
	return errhelp.FirstError(err, cerr)
}

func makeInvoker(ctx *c.Context, cfg *config.Config, errw io.Writer) (*invoker.Invoker, error) {
	l := log.New(errw, "[invoker] ", log.LstdFlags)
	v := stdflag.Verbose(ctx)

	return invoker.New(stdflag.OutDirFromCli(ctx),
		&runner.FromPlanFactory{Config: cfg.SSH},
		invoker.AllowReinvoke(ctx.Bool(flagForce)),
		invoker.ObserveCopiesWith(singleobs.Copier(l, v)...),
		invoker.ObserveMachWith(singleobs.MachNode(l, v)...),
		// The idea here is that the configuration's quantities should be lowest priority, then those in the plan,
		// and then those in the command line arguments.
		invoker.OverrideBaseQuantities(cfg.Quantities.Mach),
		invoker.OverrideQuantitiesFromPlanThen(stdflag.MachNodeQuantitySetFromCli(ctx)),
	)
}
