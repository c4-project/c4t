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

	"github.com/MattWindsor91/act-tester/internal/app/director"
	"github.com/MattWindsor91/act-tester/internal/app/mach"
	"github.com/MattWindsor91/act-tester/internal/app/query"

	"github.com/1set/gut/yos"
	"github.com/MattWindsor91/act-tester/internal/app/plan"
	"github.com/MattWindsor91/act-tester/internal/view/stdflag"

	c "github.com/urfave/cli/v2"
)

// App creates the act-tester-mkdoc app.
func App(outw, errw io.Writer) *c.App {
	return &c.App{
		Name:  "act-tester-mkdoc",
		Usage: "makes documentation for act-tester commands",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
		Writer:                 outw,
		ErrWriter:              errw,
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
	}
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

func appsToDocument(ctx *c.Context, outw io.Writer, errw io.Writer) []*c.App {
	return []*c.App{
		director.App(outw, errw),
		mach.App(outw, errw),
		plan.App(outw, errw),
		query.App(outw, errw),
		ctx.App,
	}
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
