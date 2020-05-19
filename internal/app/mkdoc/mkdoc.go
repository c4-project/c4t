// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package mkdoc contains the app definition for act-tester-mkdoc.
package mkdoc

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/MattWindsor91/act-tester/internal/app/fuzz"
	"github.com/MattWindsor91/act-tester/internal/app/lift"
	"github.com/MattWindsor91/act-tester/internal/app/rmach"

	"github.com/MattWindsor91/act-tester/internal/app/gccnt"
	"github.com/MattWindsor91/act-tester/internal/app/litmus"

	"github.com/MattWindsor91/act-tester/internal/app/analyse"
	"github.com/MattWindsor91/act-tester/internal/app/director"
	"github.com/MattWindsor91/act-tester/internal/app/mach"

	"github.com/1set/gut/yos"
	"github.com/MattWindsor91/act-tester/internal/app/plan"
	"github.com/MattWindsor91/act-tester/internal/view/stdflag"

	c "github.com/urfave/cli/v2"
)

// App creates the act-tester-mkdoc app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "act-tester-mkdoc",
		Usage: "makes documentation for act-tester commands",
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
	director.App,
	fuzz.App,
	gccnt.App,
	lift.App,
	litmus.App,
	mach.App,
	plan.App,
	analyse.App,
	rmach.App,
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
