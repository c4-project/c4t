// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/MattWindsor91/act-tester/internal/controller/mach/compiler"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/runner"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/timeout"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	c "github.com/urfave/cli/v2"

	"github.com/MattWindsor91/act-tester/internal/view/singleobs"

	"github.com/MattWindsor91/act-tester/internal/controller/mach/forward"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	bimpl "github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"
	cimpl "github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler"

	"github.com/MattWindsor91/act-tester/internal/controller/mach"

	"github.com/MattWindsor91/act-tester/internal/view"
)

const (
	defaultOutDir = "mach_results"

	flagSkipCompiler = "skip-compiler"
	flagSkipRunner   = "skip-runner"
)

func main() {
	app := c.App{
		Name:                   "act-tester-mach",
		Usage:                  "runs the machine-dependent phase of an ACT test",
		Flags:                  flags(),
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
		Action: func(ctx *c.Context) error {
			return run(ctx, os.Stdout, os.Stderr)
		},
	}
	view.LogTopError(app.Run(os.Args))
}

func flags() []c.Flag {
	return []c.Flag{
		&c.BoolFlag{
			Name:    view.FlagUseJSONLong,
			Aliases: []string{view.FlagUseJSON},
			Usage:   "emit progress reports in JSON form on stderr",
		},
		&c.BoolFlag{
			Name:  flagSkipCompiler,
			Usage: "if given, skip the compiler",
		},
		&c.BoolFlag{
			Name:  flagSkipRunner,
			Usage: "if given, skip the runner",
		},
		&c.DurationFlag{
			Name:    view.FlagCompilerTimeoutLong,
			Aliases: []string{view.FlagCompilerTimeout},
			Value:   1 * time.Minute,
			Usage:   "a `timeout` to apply to each compilation",
		},
		&c.DurationFlag{
			Name:    view.FlagRunTimeoutLong,
			Aliases: []string{view.FlagRunTimeout},
			Value:   1 * time.Minute,
			Usage:   "a `timeout` to apply to each run",
		},
		&c.IntFlag{
			Name:    view.FlagWorkerCountLong,
			Aliases: []string{view.FlagWorkerCount},
			Value:   1,
			Usage:   "number of `workers` to run in parallel",
		},
		view.OutDirCliFlag(defaultOutDir),
		view.PlanFileCliFlag(),
	}
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	cfg := makeConfig(ctx, outw, errw)
	pfile := view.PlanFileFromCli(ctx)
	return view.RunOnPlanFile(context.Background(), cfg, pfile, outw)
}

func makeConfig(ctx *c.Context, outw, errw io.Writer) *mach.Config {
	cfg := mach.Config{
		CDriver:      &cimpl.CResolve,
		RDriver:      &bimpl.BResolve,
		Stdout:       outw,
		OutDir:       view.OutDirFromCli(ctx),
		SkipCompiler: ctx.Bool(flagSkipCompiler),
		SkipRunner:   ctx.Bool(flagSkipRunner),
		Quantities:   makeQuantitySet(ctx),
	}

	setLoggerAndObservers(&cfg, errw, ctx.Bool(view.FlagUseJSONLong))
	return &cfg
}

func makeQuantitySet(ctx *c.Context) mach.QuantitySet {
	return mach.QuantitySet{
		Compiler: compiler.QuantitySet{
			Timeout: timeout.Timeout(ctx.Duration(view.FlagCompilerTimeoutLong)),
		},
		Runner: runner.QuantitySet{
			Timeout:  timeout.Timeout(ctx.Duration(view.FlagRunTimeoutLong)),
			NWorkers: ctx.Int(view.FlagWorkerCountLong),
		},
	}
}

func setLoggerAndObservers(c *mach.Config, errw io.Writer, jsonStatus bool) {
	errw = iohelp.EnsureWriter(errw)

	if jsonStatus {
		c.Logger = nil
		c.Observers = makeJsonObserver(errw)
		return
	}

	c.Logger = log.New(errw, "[mach] ", log.LstdFlags)
	c.Observers = singleobs.Builder(c.Logger)
}

func makeJsonObserver(errw io.Writer) []builder.Observer {
	return []builder.Observer{&forward.Observer{Encoder: json.NewEncoder(errw)}}
}
