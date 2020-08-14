// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package director contains the app definition for act-tester ('the director').
package director

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/act"
	"github.com/MattWindsor91/act-tester/internal/config"
	"github.com/MattWindsor91/act-tester/internal/director"
	"github.com/MattWindsor91/act-tester/internal/director/observer"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler"
	"github.com/MattWindsor91/act-tester/internal/stage/planner"
	"github.com/MattWindsor91/act-tester/internal/ux/dash"
	"github.com/mitchellh/go-homedir"

	"github.com/MattWindsor91/act-tester/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	name  = "act-tester"
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
)

// App creates the act-tester app.
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
		&c.StringFlag{
			Name:    flagMFilter,
			Aliases: []string{stdflag.FlagMachine},
			Usage:   usageMFilter,
			Value:   "",
		},
		stdflag.SubjectCyclesCliFlag(),
		stdflag.CorpusSizeCliFlag(),
		stdflag.CPUProfileCliFlag(),
	}
	return append(nflags, stdflag.ActRunnerCliFlags()...)
}

func run(ctx *c.Context, errw io.Writer) error {
	if cppath := stdflag.CPUProfileFromCli(ctx); !ystring.IsBlank(cppath) {
		stop, err := setupPprof(cppath)
		if err != nil {
			return err
		}
		defer stop()
	}

	a := stdflag.ActRunnerFromCli(ctx, errw)
	cfg, err := stdflag.ConfFileFromCli(ctx)
	if err != nil {
		return err
	}
	qs := setupQuantityOverrides(ctx)

	args := args{
		dash:    !ctx.Bool(flagNoDash),
		errw:    errw,
		mfilter: ctx.String(flagMFilter),
		files:   ctx.Args().Slice(),
	}

	return runWithArgs(cfg, qs, a, args)
}

type args struct {
	dash    bool
	errw    io.Writer
	mfilter string
	files   []string
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

func runWithArgs(cfg *config.Config, qs quantity.RootSet, a *act.Runner, args args) error {
	o, lw, err := makeObservers(cfg, args)
	if err != nil {
		return err
	}

	opts, err := makeOptions(cfg, qs, args.mfilter, lw, o...)
	if err != nil {
		_ = observer.CloseAll(o...)
		return err
	}

	e := makeEnv(a, cfg)
	d, err := director.New(e, cfg.Machines, args.files, opts...)
	if err != nil {
		_ = observer.CloseAll(o...)
		return err
	}

	// The director will close the observers.
	return d.Direct(context.Background())
}

func createResultLogFile(c *config.Config) (*os.File, error) {
	logpath, err := homedir.Expand(filepath.Join(c.OutDir, "results.log"))
	if err != nil {
		return nil, fmt.Errorf("expanding result log file path: %w", err)
	}
	logw, err := os.Create(logpath)
	if err != nil {
		return nil, fmt.Errorf("opening result log file: %w", err)
	}
	return logw, nil
}

func makeOptions(c *config.Config, qs quantity.RootSet, mfilter string, lw io.Writer, o ...observer.Observer) ([]director.Option, error) {
	glob, err := makeGlob(mfilter)
	if err != nil {
		return nil, err
	}

	l := log.New(lw, "", 0)

	opts := []director.Option{
		director.ConfigFromGlobal(c),
		director.OverrideQuantities(qs),
		director.FilterMachines(glob),
		director.ObserveWith(o...),
		director.LogWith(l),
	}
	return opts, nil
}

func makeGlob(mfilter string) (id.ID, error) {
	if ystring.IsBlank(mfilter) {
		return id.ID{}, nil
	}
	return id.TryFromString(mfilter)
}

func makeObservers(cfg *config.Config, args args) ([]observer.Observer, io.Writer, error) {
	logw, err := createResultLogFile(cfg)
	if err != nil {
		return nil, nil, err
	}
	lo, err := observer.NewLogger(logw)
	if err != nil {
		return nil, nil, err
	}
	if !args.dash {
		return []observer.Observer{lo}, args.errw, nil
	}
	do, err := dash.New()
	if err != nil {
		_ = lo.Close()
		return nil, nil, err
	}

	return []observer.Observer{do, lo}, do, nil
}

func makeEnv(a *act.Runner, c *config.Config) director.Env {
	return director.Env{
		Fuzzer:     a,
		Lifter:     &backend.BResolve,
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
				SubjectCycles: stdflag.SubjectCyclesFromCli(ctx),
			},
		},
	}
}
