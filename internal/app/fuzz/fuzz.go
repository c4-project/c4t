// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package fuzz contains the app definition for act-tester-fuzz.
package fuzz

import (
	"io"
	"log"

	"github.com/MattWindsor91/act-tester/internal/act"

	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/view/singleobs"

	"github.com/MattWindsor91/act-tester/internal/view/stdflag"

	c "github.com/urfave/cli/v2"

	"github.com/MattWindsor91/act-tester/internal/view"
)

// defaultOutDir is the default directory used for the results of the lifter.
const defaultOutDir = "fuzz_results"

// App creates the act-tester-mach app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "act-tester-fuzz",
		Usage: "runs the batch-fuzzer phase of an ACT test",
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
		stdflag.CorpusSizeCliFlag(),
		stdflag.SubjectCyclesCliFlag(),
	}
	return append(fs, stdflag.ActRunnerCliFlags()...)
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	a := stdflag.ActRunnerFromCli(ctx, errw)
	l := log.New(errw, "", 0)
	cfg := makeConfig(ctx, a, l)
	pf := stdflag.PlanFileFromCli(ctx)
	return view.RunOnPlanFile(ctx.Context, cfg, pf, outw)
}

func makeConfig(ctx *c.Context, a *act.Runner, l *log.Logger) *fuzzer.Config {
	cfg := fuzzer.Config{
		Driver:     a,
		Observers:  singleobs.Builder(l),
		Logger:     l,
		Paths:      fuzzer.NewPathset(stdflag.OutDirFromCli(ctx)),
		Quantities: *setupQuantityFlags(ctx),
	}
	return &cfg
}

func setupQuantityFlags(ctx *c.Context) *fuzzer.QuantitySet {
	return &fuzzer.QuantitySet{
		CorpusSize:    stdflag.CorpusSizeFromCli(ctx),
		SubjectCycles: stdflag.SubjectCyclesFromCli(ctx),
		NWorkers:      stdflag.WorkerCountFromCli(ctx),
	}
}
