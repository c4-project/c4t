// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package rmach contains the app definition for act-tester-rmach.
package rmach

import (
	"io"
	"log"

	"github.com/MattWindsor91/act-tester/internal/view/stdflag"

	"github.com/MattWindsor91/act-tester/internal/config"
	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/controller/rmach"

	c "github.com/urfave/cli/v2"

	"github.com/MattWindsor91/act-tester/internal/view/singleobs"

	"github.com/MattWindsor91/act-tester/internal/view"
)

const Name = "act-tester-rmach"

// App creates the act-tester-rmach app.
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
	return view.RunOnCliPlan(ctx, inv, outw)
}

func makeInvoker(ctx *c.Context, cfg *config.Config, errw io.Writer) (*rmach.Invoker, error) {
	l := log.New(errw, "[rmach] ", log.LstdFlags)
	mcfg := stdflag.MachConfigFromCli(ctx, cfg.Quantities.Mach)

	return rmach.New(stdflag.OutDirFromCli(ctx),
		stdflag.MachInvoker{
			Config: &mcfg,
		},
		rmach.ObserveWith(singleobs.RMach(l)...),
		rmach.UsePlanSSH(cfg.SSH),
	)
}
