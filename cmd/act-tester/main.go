// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"

	"github.com/MattWindsor91/act-tester/internal/pkg/director/observer"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/serviceimpl/backend"

	"github.com/MattWindsor91/act-tester/internal/pkg/ux/dash"

	"github.com/MattWindsor91/act-tester/internal/pkg/act"
	"github.com/MattWindsor91/act-tester/internal/pkg/planner"

	"github.com/MattWindsor91/act-tester/internal/pkg/config"

	"github.com/MattWindsor91/act-tester/internal/pkg/director"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

func main() {
	err := run(os.Args, os.Stderr)
	ux.LogTopError(err)
}

func run(args []string, errw io.Writer) error {
	a := act.Runner{Stderr: errw}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	ux.ActRunnerFlags(fs, &a)
	cfile := ux.ConfFileFlag(fs)
	qs := setupQuantityOverrides(fs)
	mfilter := fs.String(ux.FlagMachine, "", "A `glob` to use to filter incoming machines by ID.")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	return runWithArgs(cfile, qs, a, *mfilter, fs.Args())
}

func runWithArgs(cfile *string, qs *config.QuantitySet, a act.Runner, mfilter string, files []string) error {
	c, err := loadAndAmendConfig(cfile, qs, mfilter)
	if err != nil {
		return err
	}

	logw, err := createResultLogFile(c)
	if err != nil {
		return err
	}

	d, err := makeDirector(c, a, logw, files)
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

func loadAndAmendConfig(cfile *string, qs *config.QuantitySet, mfilter string) (*config.Config, error) {
	c, err := config.Load(*cfile)
	if err != nil {
		return nil, err
	}
	c.Quantities.Override(*qs)
	if mfilter != "" {
		if err := applyMachineFilter(mfilter, c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func makeDirector(c *config.Config, a act.Runner, logw io.Writer, files []string) (*director.Director, error) {
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

func makeDirectorConfig(c *config.Config, a act.Runner, logw io.Writer) (*director.Config, error) {
	e := makeEnv(&a, c)

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
	lo := observer.NewLogger(logw)
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
			BProbe: c,
			CProbe: c,
			SProbe: a,
		},
	}
}

func setupQuantityOverrides(fs *flag.FlagSet) *config.QuantitySet {
	var q config.QuantitySet
	// TODO(@MattWindsor91): disambiguate the corpus size argument
	ux.CorpusSizeFlag(fs, &q.Fuzz.CorpusSize)
	ux.SubjectCycleFlag(fs, &q.Fuzz.SubjectCycles)
	return &q
}
