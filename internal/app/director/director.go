// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package director contains the app definition for c4t ('the director').
package director

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/c4-project/c4t/internal/serviceimpl/backend"

	"github.com/c4-project/c4t/internal/helper/errhelp"
	"github.com/c4-project/c4t/internal/ux/directorobs"
	"golang.org/x/sync/errgroup"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/1set/gut/ystring"

	"github.com/c4-project/c4t/internal/c4f"
	"github.com/c4-project/c4t/internal/config"
	"github.com/c4-project/c4t/internal/director"
	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/serviceimpl/compiler"
	"github.com/c4-project/c4t/internal/stage/planner"

	"github.com/c4-project/c4t/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	name  = "c4t"
	usage = "runs compiler tests"

	readme = `
   This program is the main test 'director'.  It runs a series of infinite
   loops, one per machine, of the plan, fuzz, lift, invoke, and analyse test
   stages.

   It takes, as arguments, a series of Litmus files and directories that will
   serve as the test corpus across all machines and loop iterations.  As
   the plan and fuzz test stages apply subsampling, the director will
   gradually cover a range of the input, but isn't guaranteed to visit each
   input file.

   By default, the director shows its progress through an interactive terminal
   dashboard.  This dashboard can consume a large amount of resources; pass
   --` + flagNoDash + ` to disable it.

   In dashboard mode (the default), pressing Ctrl-C on the terminal stops the
   tester gracefully.  In no-dashboard mode, the tester will shut down in
   response to interrupt signals, which can usually be sent by pressing Ctrl-C
   anyway.

   Most of the director's options can be configured through the main config
   file.  Options specified on the command line, where appropriate, override
   that configuration.`

	flagMFilter  = "machine-filter"
	usageMFilter = "a `glob` to use to filter incoming machines by ID"

	flagNoDash      = "no-dashboard"
	flagNoDashShort = "D"
	usageNoDash     = "turns off the dashboard"

	flagNoFuzz      = "no-fuzz"
	flagNoFuzzShort = "F"
	usageNoFuzz     = "turns off the fuzzer stage"
)

// App creates the c4t app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:        name,
		Usage:       usage,
		Description: strings.TrimSpace(readme),
		Flags:       flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, errw)
		},
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	nflags := []c.Flag{
		stdflag.ConfFileCliFlag(),
		&c.BoolFlag{
			Name:    flagNoDash,
			Aliases: []string{flagNoDashShort},
			Usage:   usageNoDash,
		},
		&c.BoolFlag{
			Name:    flagNoFuzz,
			Aliases: []string{flagNoFuzzShort},
			Usage:   usageNoFuzz,
		},
		&c.StringFlag{
			Name:    flagMFilter,
			Aliases: []string{stdflag.FlagMachine},
			Usage:   usageMFilter,
			Value:   "",
		},
		stdflag.SubjectFuzzesCliFlag(),
		stdflag.CorpusSizeCliFlag(),
		stdflag.CPUProfileCliFlag(),
	}
	return append(nflags, stdflag.C4fRunnerCliFlags()...)
}

func run(ctx *c.Context, errw io.Writer) error {
	if cppath := stdflag.CPUProfileFromCli(ctx); !ystring.IsBlank(cppath) {
		stop, err := setupPprof(cppath)
		if err != nil {
			return err
		}
		defer stop()
	}

	a := stdflag.C4fRunnerFromCli(ctx, errw)
	cfg, err := stdflag.ConfFileFromCli(ctx)
	if err != nil {
		return err
	}
	qs := setupQuantityOverrides(ctx)

	args := args{
		dash:         !ctx.Bool(flagNoDash),
		errw:         errw,
		mfilter:      ctx.String(flagMFilter),
		files:        ctx.Args().Slice(),
		fuzzDisabled: ctx.Bool(flagNoFuzz),
	}

	return runWithArgs(ctx.Context, cfg, qs, a, args)
}

type args struct {
	dash         bool
	errw         io.Writer
	mfilter      string
	files        []string
	fuzzDisabled bool
}

func setupPprof(cppath string) (func(), error) {
	cpf, err := os.Create(cppath)
	if err != nil {
		return nil, fmt.Errorf("opening profile file: %w", err)
	}
	if err := pprof.StartCPUProfile(cpf); err != nil {
		cpf.Close()
		return nil, fmt.Errorf("starting profile to %s: %w", cppath, err)
	}
	return func() {
		pprof.StopCPUProfile()
		_ = cpf.Close()
	}, nil
}

func runWithArgs(ctx context.Context, cfg *config.Config, qs quantity.RootSet, a *c4f.Runner, args args) error {
	if err := overrideConfig(cfg, qs, args); err != nil {
		return err
	}
	o, err := directorobs.NewObs(cfg, args.dash)
	if err != nil {
		return err
	}
	err = runWithObs(ctx, cfg, args, a, o)
	cerr := o.Close()
	return errhelp.FirstError(err, cerr)
}

func runWithObs(ctx context.Context, cfg *config.Config, args args, a *c4f.Runner, o *directorobs.Obs) error {
	glob, err := makeGlob(args.mfilter)
	if err != nil {
		return err
	}
	d, err := makeDirector(cfg, glob, a, o)
	if err != nil {
		return err
	}

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// TODO(@MattWindsor91): is this nesting of errgroups inefficient?
	eg, ectx := errgroup.WithContext(cctx)
	eg.Go(func() error {
		return d.Direct(ectx)
	})
	eg.Go(func() error {
		return o.Run(ectx, cancel)
	})
	return eg.Wait()
}

func makeDirector(cfg *config.Config, glob id.ID, a *c4f.Runner, obs *directorobs.Obs) (*director.Director, error) {
	return director.New(makeEnv(a, cfg), cfg.Machines, cfg.Paths.Inputs,
		director.ConfigFromGlobal(cfg),
		director.FilterMachines(glob),
		director.ObserveWith(obs.Observers()...),
	)
}

func overrideConfig(cfg *config.Config, qs quantity.RootSet, args args) error {
	cfg.OverrideQuantities(qs)
	if args.fuzzDisabled {
		cfg.DisableFuzz()
	}
	return cfg.OverrideInputs(args.files)
}

func makeGlob(mfilter string) (id.ID, error) {
	if ystring.IsBlank(mfilter) {
		return id.ID{}, nil
	}
	return id.TryFromString(mfilter)
}

func makeEnv(a *c4f.Runner, c *config.Config) director.Env {
	return director.Env{
		Fuzzer:     a,
		BResolver:  &backend.Resolve,
		CInspector: &compiler.CResolve,
		Planner: planner.Source{
			BProbe:  c,
			CLister: c.Machines,
			SProbe:  a,
		},
	}
}

func setupQuantityOverrides(ctx *c.Context) quantity.RootSet {
	// TODO(@MattWindsor91): disambiguate the corpus size argument
	return quantity.RootSet{
		MachineSet: quantity.MachineSet{
			Fuzz: quantity.FuzzSet{
				CorpusSize:    stdflag.CorpusSizeFromCli(ctx),
				SubjectCycles: stdflag.SubjectFuzzesFromCli(ctx),
			},
		},
	}
}
