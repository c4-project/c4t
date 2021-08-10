// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package mkdoc contains the app definition for c4t-mkdoc.
package mkdoc

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/c4-project/c4t/internal/app/mkdb"

	"github.com/c4-project/c4t/internal/app/stat"

	"github.com/c4-project/c4t/internal/app/config"

	"github.com/c4-project/c4t/internal/app/backend"
	"github.com/c4-project/c4t/internal/app/obs"

	"github.com/c4-project/c4t/internal/app/coverage"

	"github.com/c4-project/c4t/internal/app/analyse"
	"github.com/c4-project/c4t/internal/app/invoke"
	"github.com/c4-project/c4t/internal/app/perturb"
	"github.com/c4-project/c4t/internal/app/setc"

	"github.com/c4-project/c4t/internal/app/fuzz"
	"github.com/c4-project/c4t/internal/app/lift"

	"github.com/c4-project/c4t/internal/app/gccnt"

	"github.com/c4-project/c4t/internal/app/director"
	"github.com/c4-project/c4t/internal/app/mach"

	"github.com/1set/gut/yos"
	"github.com/c4-project/c4t/internal/app/plan"
	"github.com/c4-project/c4t/internal/ux/stdflag"

	c "github.com/urfave/cli/v2"
)

// App creates the c4t-mkdoc app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "c4t-mkdoc",
		Usage: "makes documentation for c4t commands",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	return []c.Flag{
		stdflag.OutDirCliFlag("docs"),
	}
}

func run(ctx *c.Context, outw io.Writer, errw io.Writer) error {
	outdir := stdflag.OutDirFromCli(ctx)
	for _, app := range appsToDocument(ctx, outw, errw) {
		if err := runApp(outdir, app); err != nil {
			return fmt.Errorf("in app %s: %w", app.Name, err)
		}
	}
	return nil
}

var appFuncs = [...]func(io.Writer, io.Writer) *c.App{
	analyse.App,
	backend.App,
	config.App,
	coverage.App,
	director.App,
	fuzz.App,
	gccnt.App,
	invoke.App,
	lift.App,
	mach.App,
	mkdb.App,
	obs.App,
	perturb.App,
	plan.App,
	setc.App,
	stat.App,
}

func appsToDocument(ctx *c.Context, outw io.Writer, errw io.Writer) []*c.App {
	apps := make([]*c.App, len(appFuncs)+1)
	apps[0] = ctx.App
	for i, f := range appFuncs {
		apps[i+1] = f(outw, errw)
	}
	return apps
}

func runApp(outroot string, app *c.App) error {
	name := app.Name
	outdir := filepath.Join(outroot, name)
	if err := yos.MakeDir(outdir); err != nil {
		return fmt.Errorf("making dir for %s: %w", name, err)
	}

	for mname, m := range methodsOf(app) {
		if err := m.run(outdir); err != nil {
			return fmt.Errorf("making %s for %s: %w", mname, name, err)
		}
	}
	return nil
}
