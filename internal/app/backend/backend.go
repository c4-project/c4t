// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/c4-project/c4t/internal/c4f"

	"github.com/c4-project/c4t/internal/helper/errhelp"

	"github.com/c4-project/c4t/internal/model/recipe"

	"github.com/c4-project/c4t/internal/model/service"

	"github.com/c4-project/c4t/internal/helper/srvrun"
	"github.com/c4-project/c4t/internal/subject/obs"

	"github.com/c4-project/c4t/internal/config"
	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/model/service/backend"
	backend2 "github.com/c4-project/c4t/internal/serviceimpl/backend"
	"github.com/c4-project/c4t/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	// Name is the name of the backend binary.
	Name  = "c4t-backend"
	usage = "runs backends standalone"

	readme = `
    This program runs lifting backends directly, in their standalone mode, and parses their results into the usual C4
    observation format.

    Doing so avoids the need to produce a testing plan, but gives less control over any compilation or other
    intermediate actions.

    The backend to run may be controlled by the -` + flagBackendIDGlob + ` glob ID, which filters on the user-defined
    name of the backend, and the -` + flagBackendStyleGlob + ` glob ID, which filters on the style of the backend.  The
    first configured backend satisfying all of the given constraints is used.`

	flagBackendIDGlob         = "backend-id"
	flagBackendIDGlobShort    = "n"
	usageBackendIDGlob        = "filter to backends whose names match `GLOB`"
	flagBackendStyleGlob      = "backend-style"
	flagBackendStyleGlobShort = "s"
	usageBackendStyleGlob     = "filter to backends whose styles match `GLOB`"
	flagArchID                = "arch"
	flagArchIDShort           = "a"
	usageArchID               = "ID of `ARCH` to target for architecture-dependent backends"
	flagDryRun                = "dry-run"
	flagDryRunShort           = "d"
	usageDryRun               = "if true, print any external commands run instead of running them"
	flagTimeout               = "timeout"
	flagTimeoutShort          = "t"
	usageTimeout              = "`DURATION` to wait before trying to stop the backend"
	flagGrace                 = "grace"
	flagGraceShort            = "g"
	usageGrace                = "`DURATION` to wait between sigterm and sigkill when timing out"
)

// App is the c4-backend app.
func App(outw, errw io.Writer) *c.App {
	a := &c.App{
		Name:        Name,
		Usage:       usage,
		Description: readme,
		Flags:       flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetPlanAppSettings(a, outw, errw)
}

func flags() []c.Flag {
	ownFlags := []c.Flag{
		stdflag.ConfFileCliFlag(),
		&c.GenericFlag{Name: flagArchID, Aliases: []string{flagArchIDShort}, Usage: usageArchID, Value: &id.ID{}},
		&c.GenericFlag{Name: flagBackendIDGlob, Aliases: []string{flagBackendIDGlobShort}, Usage: usageBackendIDGlob, Value: &id.ID{}},
		&c.GenericFlag{Name: flagBackendStyleGlob, Aliases: []string{flagBackendStyleGlobShort}, Usage: usageBackendStyleGlob, Value: &id.ID{}},
		&c.BoolFlag{Name: flagDryRun, Aliases: []string{flagDryRunShort}, Usage: usageDryRun},
		&c.DurationFlag{Name: flagTimeout, Aliases: []string{flagTimeoutShort}, Usage: usageTimeout},
		&c.DurationFlag{Name: flagGrace, Aliases: []string{flagGraceShort}, Usage: usageGrace},
	}
	return append(ownFlags, stdflag.C4fRunnerCliFlags()...)
}

func run(ctx *c.Context, outw io.Writer, errw io.Writer) error {
	cfg, err := stdflag.ConfFileFromCli(ctx)
	if err != nil {
		return fmt.Errorf("while getting config: %w", err)
	}
	crun := stdflag.C4fRunnerFromCli(ctx, errw)

	td, err := ioutil.TempDir("", "c4t-backend")
	if err != nil {
		return err
	}

	fn, err := inputNameFromCli(ctx)
	if err != nil {
		return err
	}

	b, err := getBackend(cfg, criteriaFromCli(ctx))
	if err != nil {
		return err
	}

	j, err := jobFromCli(ctx, fn, crun, td)
	if err != nil {
		return err
	}
	perr := runParseAndDump(ctx, outw, b, j, makeRunner(ctx, errw))
	derr := os.RemoveAll(td)
	return errhelp.FirstError(perr, derr)
}

func jobFromCli(ctx *c.Context, fn string, c4f *c4f.Runner, td string) (backend.LiftJob, error) {
	in, err := backend.InputFromFile(ctx.Context, fn, c4f)
	if err != nil {
		return backend.LiftJob{}, err
	}

	j := backend.LiftJob{
		Arch: idFromCli(ctx, flagArchID),
		In:   in,
		Out:  backend.LiftOutput{Dir: td, Target: backend.ToStandalone},
	}
	return j, nil
}

func makeRunner(ctx *c.Context, errw io.Writer) service.Runner {
	// TODO(@MattWindsor91): the backend logic isn't very resilient against having external commands not run.
	if ctx.Bool(flagDryRun) {
		return srvrun.DryRunner{Writer: errw}
	}
	// TODO(@MattWindsor91): use grace in the rest of c4t
	return srvrun.NewExecRunner(srvrun.StderrTo(errw), srvrun.WithGrace(ctx.Duration(flagGrace)))
}

func runParseAndDump(ctx *c.Context, outw io.Writer, b backend.Backend, j backend.LiftJob, xr service.Runner) error {
	var o obs.Obs

	to := ctx.Duration(flagTimeout)
	if err := runAndParse(ctx.Context, to, b, j, &o, xr); err != nil {
		return err
	}

	e := json.NewEncoder(outw)
	e.SetIndent("", "\t")
	return e.Encode(o)
}

func runAndParse(ctx context.Context, to time.Duration, b backend.Backend, j backend.LiftJob, o *obs.Obs, xr service.Runner) error {
	// TODO(@MattWindsor91): clean this function up, eg making a separate struct...

	r, err := liftWithTimeout(ctx, to, b, j, xr)
	if err != nil {
		return err
	}

	if r.Output != recipe.OutNothing {
		return fmt.Errorf("can't handle recipes with outputs: %s", r.Output)
	}

	for _, fname := range r.Paths() {
		if err := parseFile(ctx, b, o, fname); err != nil {
			return err
		}
	}
	return nil
}

func liftWithTimeout(ctx context.Context, to time.Duration, b backend.Backend, j backend.LiftJob, xr service.Runner) (recipe.Recipe, error) {
	cf := func() {}
	if to != 0 {
		ctx, cf = context.WithTimeout(ctx, to)
	}
	defer cf()

	// TODO(@MattWindsor91): deduplicate with runAndParseBin?.
	return b.Lift(ctx, j, xr)
}

func parseFile(ctx context.Context, b backend.Backend, o *obs.Obs, fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("can't open output file %s: %w", fname, err)
	}
	perr := b.ParseObs(ctx, f, o)
	cerr := f.Close()
	return errhelp.FirstError(perr, cerr)
}

func inputNameFromCli(ctx *c.Context) (string, error) {
	if ctx.Args().Len() != 1 {
		return "", errors.New("expected one argument")
	}
	return ctx.Args().First(), nil
}

func getBackend(cfg *config.Config, c backend.Criteria) (backend.Backend, error) {
	spec, err := cfg.FindBackend(c)
	if err != nil {
		return nil, fmt.Errorf("while finding backend: %w", err)
	}

	b, err := backend2.Resolve.Resolve(spec.Spec)
	if err != nil {
		return nil, fmt.Errorf("while resolving backend %s: %w", spec.ID, err)
	}
	return b, nil
}

func criteriaFromCli(ctx *c.Context) backend.Criteria {
	return backend.Criteria{
		IDGlob:    idFromCli(ctx, flagBackendIDGlob),
		StyleGlob: idFromCli(ctx, flagBackendStyleGlob),
	}
}

func idFromCli(ctx *c.Context, flag string) id.ID {
	return *(ctx.Generic(flag).(*id.ID))
}
