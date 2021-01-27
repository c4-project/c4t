// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package fuzz contains the app definition for c4t-fuzz.
package fuzz

import (
	"io"
	"log"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/config"

	"github.com/c4-project/c4t/internal/stage/fuzzer"

	"github.com/c4-project/c4t/internal/ux/singleobs"

	"github.com/c4-project/c4t/internal/ux/stdflag"

	c "github.com/urfave/cli/v2"

	"github.com/c4-project/c4t/internal/ux"
)

// defaultOutDir is the default directory used for the results of the fuzzer.
const defaultOutDir = "fuzz_results"

// App creates the c4t-fuzz app.
func App(outw, errw io.Writer) *c.App {
	a := &c.App{
		Name:  "c4t-fuzz",
		Usage: "runs the batch-fuzzer phase of a C4 test",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetPlanAppSettings(a, outw, errw)
}

func flags() []c.Flag {
	fs := []c.Flag{
		stdflag.VerboseFlag(),
		stdflag.ConfFileCliFlag(),
		stdflag.OutDirCliFlag(defaultOutDir),
		stdflag.CorpusSizeCliFlag(),
		stdflag.SubjectCyclesCliFlag(),
	}
	return append(fs, stdflag.C4fRunnerCliFlags()...)
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	cfg, err := stdflag.ConfFileFromCli(ctx)
	if err != nil {
		return err
	}

	a := stdflag.C4fRunnerFromCli(ctx, errw)
	l := log.New(errw, "", 0)
	f, err := makeFuzzer(ctx, cfg, a, l)
	if err != nil {
		return err
	}
	return ux.RunOnCliPlan(ctx, f, outw)
}

func makeFuzzer(ctx *c.Context, cfg *config.Config, drv fuzzer.Driver, l *log.Logger) (*fuzzer.Fuzzer, error) {
	return fuzzer.New(
		drv,
		fuzzer.NewPathset(stdflag.OutDirFromCli(ctx)),
		fuzzer.ObserveWith(singleobs.Builder(l, stdflag.Verbose(ctx))...),
		fuzzer.OverrideQuantities(cfg.Quantities.Fuzz),
		fuzzer.OverrideQuantities(setupQuantityFlags(ctx)),
		fuzzer.UseConfig(cfg.Fuzz),
	)
}

func setupQuantityFlags(ctx *c.Context) quantity.FuzzSet {
	return quantity.FuzzSet{
		CorpusSize:    stdflag.CorpusSizeFromCli(ctx),
		SubjectCycles: stdflag.SubjectCyclesFromCli(ctx),
		NWorkers:      stdflag.WorkerCountFromCli(ctx),
	}
}
