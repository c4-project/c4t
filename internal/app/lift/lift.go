// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package lift contains the app definition for c4t-lift.
package lift

import (
	"io"
	"log"

	"github.com/c4-project/c4t/internal/serviceimpl/backend/resolver"

	"github.com/c4-project/c4t/internal/stage/lifter"

	"github.com/c4-project/c4t/internal/ux/singleobs"

	"github.com/c4-project/c4t/internal/ux/stdflag"

	c "github.com/urfave/cli/v2"

	"github.com/c4-project/c4t/internal/ux"
)

// defaultOutDir is the default directory used for the results of the lifter.
const defaultOutDir = "lift_results"

// App creates the c4t-lift app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "c4t-lift",
		Usage: "runs the lifter phase of an ACT test",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetPlanAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	fs := []c.Flag{
		stdflag.VerboseFlag(),
		stdflag.OutDirCliFlag(defaultOutDir),
	}
	return append(fs, stdflag.ActRunnerCliFlags()...)
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	l := log.New(errw, "", 0)
	lft, err := makeLifter(ctx, l, errw)
	if err != nil {
		return err
	}
	pf, err := stdflag.PlanFileFromCli(ctx)
	if err != nil {
		return err
	}
	return ux.RunOnPlanFile(ctx.Context, lft, pf, outw)
}

func makeLifter(ctx *c.Context, l *log.Logger, errw io.Writer) (*lifter.Lifter, error) {
	return lifter.New(
		&resolver.Resolve,
		lifter.NewPathset(stdflag.OutDirFromCli(ctx)),
		lifter.ObserveWith(singleobs.Builder(l, stdflag.Verbose(ctx))...),
		lifter.SendStderrTo(errw),
	)
}
