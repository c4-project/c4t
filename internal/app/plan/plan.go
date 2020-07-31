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

	"github.com/1set/gut/yos"

	"github.com/MattWindsor91/act-tester/internal/act"

	"github.com/MattWindsor91/act-tester/internal/config"
	"github.com/MattWindsor91/act-tester/internal/plan"
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
			return run(ctx, os.Stderr)
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
		stdflag.OutDirCliFlag(""),
	}
	return append(ownFlags, stdflag.ActRunnerCliFlags()...)
}

func run(ctx *c.Context, errw io.Writer) error {
	cfg, err := stdflag.ConfFileFromCli(ctx)
	if err != nil {
		return err
	}

	pr, err := makePlanner(ctx, cfg, errw)
	if err != nil {
		return err
	}

	ps, err := pr.Plan(ctx.Context, cfg.Machines, ctx.Args().Slice()...)
	if err != nil {
		return err
	}

	return writePlans(stdflag.OutDirFromCli(ctx), ps)
}

func writePlans(outdir string, ps map[string]plan.Plan) error {
	if err := yos.MakeDir(outdir); err != nil {
		return err
	}
	for n, p := range ps {
		file := fmt.Sprintf("plan.%s.json", n)
		if err := p.WriteFile(yos.JoinPath(outdir, file), plan.WriteHuman); err != nil {
			return err
		}
	}
	return nil
}

func makePlanner(ctx *c.Context, cfg *config.Config, errw io.Writer) (*planner.Planner, error) {
	a := stdflag.ActRunnerFromCli(ctx, errw)

	qs := quantities(ctx)
	src := source(a, cfg)

	l := log.New(errw, "[planner] ", log.LstdFlags)

	return planner.New(
		src,
		planner.ObserveWith(singleobs.Planner(l)...),
		planner.OverrideQuantities(qs),
		planner.FilterCompilers(ctx.String(flagCompilerFilter)),
	)
}

func source(a *act.Runner, cfg *config.Config) planner.Source {
	return planner.Source{
		BProbe:  cfg,
		CLister: cfg.Machines,
		SProbe:  a,
	}
}

func quantities(ctx *c.Context) planner.QuantitySet {
	return planner.QuantitySet{
		NWorkers: stdflag.WorkerCountFromCli(ctx),
	}
}
