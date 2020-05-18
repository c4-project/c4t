// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package lift contains the app definition for act-tester-lift.
package lift

import (
	"io"
	"log"

	"github.com/MattWindsor91/act-tester/internal/controller/lifter"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"

	"github.com/MattWindsor91/act-tester/internal/view/singleobs"

	"github.com/MattWindsor91/act-tester/internal/view/stdflag"

	c "github.com/urfave/cli/v2"

	"github.com/MattWindsor91/act-tester/internal/view"
)

// defaultOutDir is the default directory used for the results of the lifter.
const defaultOutDir = "lift_results"

// App creates the act-tester-lift app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "act-tester-lift",
		Usage: "runs the harness-lifter phase of an ACT test",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	fs := []c.Flag{
		stdflag.OutDirCliFlag(defaultOutDir),
		stdflag.PlanFileCliFlag(),
	}
	return append(fs, stdflag.ActRunnerCliFlags()...)
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	l := log.New(errw, "", 0)
	cfg := makeConfig(ctx, l, errw)
	pf := stdflag.PlanFileFromCli(ctx)
	return view.RunOnPlanFile(ctx.Context, cfg, pf, outw)
}

func makeConfig(ctx *c.Context, l *log.Logger, errw io.Writer) *lifter.Config {
	return &lifter.Config{
		Maker:     &backend.BResolve,
		Logger:    l,
		Observers: singleobs.Builder(l),
		Paths:     lifter.NewPathset(stdflag.OutDirFromCli(ctx)),
		Stderr:    errw,
	}
}
