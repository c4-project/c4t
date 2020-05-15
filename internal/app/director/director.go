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
	return &c.App{
		Name:  "act-tester",
		Usage: "makes documentation for act-tester commands",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, errw)
		},
		Writer:                 outw,
		ErrWriter:              errw,
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
	}
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

func runWithArgs(cfg *config.Config, qs *config.QuantitySet, a *act.Runner, mfilter string, files []string) error {
	if err := amendConfig(cfg, qs, mfilter); err != nil {
		return err
	}

	logw, err := createResultLogFile(cfg)
	if err != nil {
		return err
	}

	d, err := makeDirector(cfg, a, logw, files)
	if err != nil {
		_ = logw.Close()
		return err
	}

	if derr := d.Direct(context.Background()); derr != nil {
		_ = logw.Close()
		return derr
	}

	return logw.Close()
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

func amendConfig(cfg *config.Config, qs *config.QuantitySet, mfilter string) error {
	cfg.Quantities.Override(*qs)
	if mfilter != "" {
		if err := applyMachineFilter(mfilter, cfg); err != nil {
			return err
		}
	}
	return nil
}

func makeDirector(c *config.Config, a *act.Runner, logw io.Writer, files []string) (*director.Director, error) {
	dc, err := makeDirectorConfig(c, a, logw)
	if err != nil {
		return nil, err
	}

	d, err := director.New(dc, files)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func makeDirectorConfig(c *config.Config, a *act.Runner, logw io.Writer) (*director.Config, error) {
	e := makeEnv(a, c)

	mids, err := c.MachineIDs()
	if err != nil {
		return nil, err
	}

	o, lw, err := makeObservers(mids, logw)
	if err != nil {
		return nil, err
	}

	l := log.New(lw, "", 0)

	dc, err := director.ConfigFromGlobal(c, l, e, o)
	if err != nil {
		return nil, err
	}
	return dc, nil
}

func makeObservers(mids []id.ID, logw io.Writer) ([]observer.Observer, io.Writer, error) {
	do, err := dash.New(mids)
	if err != nil {
		return nil, nil, err
	}
	lo, err := observer.NewLogger(logw)
	return []observer.Observer{do, lo}, do, err
}

func applyMachineFilter(mfilter string, c *config.Config) error {
	mglob, err := id.TryFromString(mfilter)
	if err != nil {
		return fmt.Errorf("parsing machine filter: %w", err)
	}
	if err := c.FilterMachines(mglob); err != nil {
		return fmt.Errorf("applying machine filter: %w", err)
	}
	return nil
}

func makeEnv(a *act.Runner, c *config.Config) director.Env {
	return director.Env{
		Fuzzer: a,
		Lifter: &backend.BResolve,
		Planner: planner.Source{
			BProbe:     c,
			CLister:    c,
			CInspector: &compiler.CResolve,
			SProbe:     a,
		},
	}
}

func setupQuantityOverrides(ctx *c.Context) *config.QuantitySet {
	// TODO(@MattWindsor91): disambiguate the corpus size argument
	return &config.QuantitySet{
		Fuzz: fuzzer.QuantitySet{
			CorpusSize:    stdflag.CorpusSizeFromCli(ctx),
			SubjectCycles: stdflag.SubjectCyclesFromCli(ctx),
		},
	}
}
