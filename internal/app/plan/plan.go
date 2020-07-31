// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package plan contains the app definition for act-tester-plan.
package plan

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/MattWindsor91/act-tester/internal/act"

	"github.com/MattWindsor91/act-tester/internal/model/machine"

	"github.com/MattWindsor91/act-tester/internal/config"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler"
	"github.com/MattWindsor91/act-tester/internal/stage/planner"
	"github.com/MattWindsor91/act-tester/internal/ux/singleobs"
	"github.com/MattWindsor91/act-tester/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	usageMach = "ID of machine to use for this test plan"

	flagCompilerFilter  = "filter-compiler"
	usageCompilerFilter = "`glob` to use to filter compilers to enable"
)

// App creates the act-tester-plan app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "act-tester-plan",
		Usage: "runs the planning phase of an ACT test standalone",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, os.Stdout, os.Stderr)
		},
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	ownFlags := []c.Flag{
		stdflag.ConfFileCliFlag(),
		&c.StringFlag{
			Name:  stdflag.FlagMachine,
			Usage: usageMach,
		},
		&c.StringFlag{
			Name:  flagCompilerFilter,
			Usage: usageCompilerFilter,
		},
		stdflag.WorkerCountCliFlag(),
	}
	return append(ownFlags, stdflag.ActRunnerCliFlags()...)
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	pr, err := makePlanner(ctx, errw)
	if err != nil {
		return err
	}

	p, err := pr.Plan(ctx.Context)
	if err != nil {
		return err
	}

	return p.Write(outw, plan.WriteHuman)
}

func makePlanner(ctx *c.Context, errw io.Writer) (*planner.Planner, error) {
	a := stdflag.ActRunnerFromCli(ctx, errw)

	cfg, err := stdflag.ConfFileFromCli(ctx)
	if err != nil {
		return nil, err
	}

	qs := quantities(ctx)
	src := source(a, cfg)
	fs := ctx.Args().Slice()

	midstr := ctx.String(stdflag.FlagMachine)
	mach, err := getMachine(cfg, midstr)
	if err != nil {
		return nil, err
	}

	l := log.New(errw, "", 0)

	return planner.New(
		src,
		mach,
		fs,
		planner.LogWith(l),
		planner.ObserveWith(singleobs.Planner(l)...),
		planner.OverrideQuantities(qs),
		planner.FilterCompilers(ctx.String(flagCompilerFilter)),
	)
}

func source(a *act.Runner, cfg *config.Config) planner.Source {
	return planner.Source{
		BProbe:     cfg,
		CLister:    cfg.Machines,
		CInspector: &compiler.CResolve,
		SProbe:     a,
	}
}

func quantities(ctx *c.Context) planner.QuantitySet {
	return planner.QuantitySet{
		NWorkers: stdflag.WorkerCountFromCli(ctx),
	}
}

func getMachine(cfg *config.Config, midstr string) (machine.Named, error) {
	mid, err := id.TryFromString(midstr)
	if err != nil {
		return machine.Named{}, err
	}

	mach, ok := cfg.Machines[midstr]
	if !ok {
		return machine.Named{}, fmt.Errorf("no such machine: %s", midstr)
	}
	m := machine.Named{
		ID:      mid,
		Machine: mach.Machine,
	}
	return m, nil
}
