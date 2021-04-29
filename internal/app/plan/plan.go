// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package plan contains the app definition for c4t-plan.
package plan

import (
	"fmt"
	"io"
	"log"
	"os"

	backend2 "github.com/c4-project/c4t/internal/serviceimpl/backend"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/id"

	"github.com/1set/gut/ystring"
	"github.com/c4-project/c4t/internal/machine"

	"github.com/1set/gut/yos"

	"github.com/c4-project/c4t/internal/config"
	"github.com/c4-project/c4t/internal/plan"
	"github.com/c4-project/c4t/internal/stage/planner"
	"github.com/c4-project/c4t/internal/ux/singleobs"
	"github.com/c4-project/c4t/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	flagCompilerFilter  = "filter-compilers"
	usageCompilerFilter = "`glob` to use to filter compilers to enable"

	flagMachineFilter  = "filter-machines"
	usageMachineFilter = "`glob` to use to filter machines to plan"
)

// App creates the c4t-plan app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "c4t-plan",
		Usage: "runs the planning phase of a C4 test standalone",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, os.Stdout, os.Stderr)
		},
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	ownFlags := []c.Flag{
		stdflag.VerboseFlag(),
		stdflag.ConfFileCliFlag(),
		&c.StringFlag{
			Name:        flagCompilerFilter,
			Aliases:     []string{stdflag.FlagCompiler},
			Usage:       usageCompilerFilter,
			DefaultText: "all compilers",
		},
		&c.StringFlag{
			Name:        flagMachineFilter,
			Aliases:     []string{stdflag.FlagMachine},
			Usage:       usageMachineFilter,
			DefaultText: "all machines",
		},
		stdflag.WorkerCountCliFlag(),
		stdflag.OutDirCliFlag(""),
	}
	return append(ownFlags, stdflag.C4fRunnerCliFlags()...)
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	cfg, err := stdflag.ConfigFromCli(ctx)
	if err != nil {
		return err
	}

	pr, err := makePlanner(ctx, cfg, errw)
	if err != nil {
		return err
	}

	ms, err := machines(ctx, cfg)
	if err != nil {
		return err
	}
	dir, err := outDir(ctx, ms)
	if err != nil {
		return err
	}

	ps, err := pr.Plan(ctx.Context, ms, ctx.Args().Slice()...)
	if err != nil {
		return err
	}

	return writePlans(outw, dir, ps)
}

func machines(ctx *c.Context, cfg *config.Config) (machine.ConfigMap, error) {
	// TODO(@MattWindsor91): maybe merge Filter into Machines.
	ms, err := cfg.Machines()
	if err != nil {
		return nil, err
	}
	midstr := ctx.String(flagMachineFilter)
	if ystring.IsBlank(midstr) {
		return ms, nil
	}
	return globbedMachines(midstr, ms)
}

func globbedMachines(midstr string, configMap machine.ConfigMap) (machine.ConfigMap, error) {
	mid, err := id.TryFromString(midstr)
	if err != nil {
		return nil, err
	}
	return configMap.Filter(mid)
}

func outDir(ctx *c.Context, ms machine.ConfigMap) (string, error) {
	dir := stdflag.OutDirFromCli(ctx)
	if ystring.IsBlank(dir) && len(ms) != 1 {
		return "", fmt.Errorf("must specify directory if planning multiple machines (have %d)", len(ms))
	}
	return dir, nil
}

func writePlans(outw io.Writer, outdir string, ps plan.Map) error {
	// Assuming that outDir above has dealt with the case whereby there is no output directory but multiple plans.
	if ystring.IsBlank(outdir) {
		return writePlansToWriter(outw, ps)
	}
	return writePlansToDir(outdir, ps)
}

func writePlansToWriter(w io.Writer, ps plan.Map) error {
	for _, p := range ps {
		if err := p.Write(w, plan.WriteHuman); err != nil {
			return err
		}
	}
	return nil
}

func writePlansToDir(outdir string, ps plan.Map) error {
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
	a := stdflag.C4fRunnerFromCli(ctx, errw)

	qs := quantities(ctx)
	src := source(a, cfg)

	l := log.New(errw, "[planner] ", log.LstdFlags)

	return planner.New(
		src,
		planner.ObserveWith(singleobs.Planner(l, stdflag.Verbose(ctx))...),
		planner.OverrideQuantities(qs),
		planner.FilterCompilers(ctx.String(flagCompilerFilter)),
	)
}

func source(a planner.SubjectProber, cfg *config.Config) planner.Source {
	// TODO(@MattWindsor91): clean all of this up across plan, director, etc.
	return planner.Source{
		BProbe: config.BackendFinder{Config: cfg, Resolver: &backend2.Resolve},
		SProbe: a,
	}
}

func quantities(ctx *c.Context) quantity.PlanSet {
	return quantity.PlanSet{
		NWorkers: stdflag.WorkerCountFromCli(ctx),
	}
}
