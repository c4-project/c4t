// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package invoker contains the app definition for act-tester-invoke.
package invoke

import (
	"io"
	"log"

	"github.com/MattWindsor91/act-tester/internal/helper/errhelp"

	"github.com/MattWindsor91/act-tester/internal/ux/stdflag"

	"github.com/MattWindsor91/act-tester/internal/config"
	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/stage/invoker"

	c "github.com/urfave/cli/v2"

	"github.com/MattWindsor91/act-tester/internal/ux/singleobs"

	"github.com/MattWindsor91/act-tester/internal/ux"
)

const Name = "act-tester-invoke"

// App creates the act-tester-invoke app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  Name,
		Usage: "runs the machine-dependent phase of an ACT test, potentially remotely",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetPlanAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	ownFlags := []c.Flag{
		stdflag.ConfFileCliFlag(),
	}
	return append(ownFlags, stdflag.MachCliFlags()...)
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	cfg, err := stdflag.ConfFileFromCli(ctx)
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
	mcfg := stdflag.MachConfigFromCli(ctx, cfg.Quantities.Mach)

	return invoker.New(stdflag.OutDirFromCli(ctx),
		stdflag.MachInvoker{
			Config: &mcfg,
		},
		invoker.ObserveWith(singleobs.Invoker(l)...),
		invoker.UsePlanSSH(cfg.SSH),
	)
}
