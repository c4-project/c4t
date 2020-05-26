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

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/act"
	"github.com/MattWindsor91/act-tester/internal/config"
	"github.com/MattWindsor91/act-tester/internal/controller/planner"
	"github.com/MattWindsor91/act-tester/internal/director"
	"github.com/MattWindsor91/act-tester/internal/director/observer"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler"
	"github.com/MattWindsor91/act-tester/internal/view/dash"
	"github.com/mitchellh/go-homedir"

	"github.com/MattWindsor91/act-tester/internal/view/stdflag"
	c "github.com/urfave/cli/v2"
)

// App creates the act-tester app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "act-tester",
		Usage: "makes documentation for act-tester commands",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, errw)
		},
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

const flagMFilter = "machine-filter"

func flags() []c.Flag {
	nflags := []c.Flag{
		stdflag.ConfFileCliFlag(),
		&c.StringFlag{
			Name:    flagMFilter,
			Aliases: []string{stdflag.FlagMachine},
			Usage:   "A `glob` to use to filter incoming machines by ID.",
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
	mfilter := ctx.String(flagMFilter)

	return runWithArgs(cfg, qs, a, mfilter, ctx.Args().Slice())
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

func runWithArgs(cfg *config.Config, qs config.QuantitySet, a *act.Runner, mfilter string, files []string) error {
	o, lw, err := makeObservers(cfg)
	if err != nil {
		return err
	}

	opts, err := makeOptions(cfg, qs, mfilter, lw, o...)
	if err != nil {
		_ = observer.CloseAll(o...)
		return err
	}

	e := makeEnv(a, cfg)
	d, err := director.New(e, cfg.Machines, files, opts...)
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

func makeOptions(c *config.Config, qs config.QuantitySet, mfilter string, lw io.Writer, o ...observer.Observer) ([]director.Option, error) {
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

func makeObservers(cfg *config.Config) ([]observer.Observer, io.Writer, error) {
	logw, err := createResultLogFile(cfg)
	if err != nil {
		return nil, nil, err
	}

	do, err := dash.New()
	if err != nil {
		return nil, nil, err
	}
	lo, err := observer.NewLogger(logw)
	if err != nil {
		_ = do.Close()
		return nil, nil, err
	}

	return []observer.Observer{do, lo}, do, nil
}

func makeEnv(a *act.Runner, c *config.Config) director.Env {
	return director.Env{
		Fuzzer: a,
		Lifter: &backend.BResolve,
		Planner: planner.Source{
			BProbe:     c,
			CLister:    c.Machines,
			CInspector: &compiler.CResolve,
			SProbe:     a,
		},
	}
}

func setupQuantityOverrides(ctx *c.Context) config.QuantitySet {
	// TODO(@MattWindsor91): disambiguate the corpus size argument
	return config.QuantitySet{
		Fuzz: fuzzer.QuantitySet{
			CorpusSize:    stdflag.CorpusSizeFromCli(ctx),
			SubjectCycles: stdflag.SubjectCyclesFromCli(ctx),
		},
	}
}
