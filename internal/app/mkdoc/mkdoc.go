// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package mkdoc contains the app definition for act-tester-mkdoc.
package mkdoc

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/1set/gut/yos"
	"github.com/MattWindsor91/act-tester/internal/app/plan"
	"github.com/MattWindsor91/act-tester/internal/view/stdflag"

	c "github.com/urfave/cli/v2"
)

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
	for _, app := range []*c.App{
		plan.App(outw, errw),
		ctx.App,
	} {
		if err := runApp(outdir, app); err != nil {
			return fmt.Errorf("in app %s: %w", app.Name, err)
		}
	}
	return nil
}

type method struct {
	name string
	make func() (string, error)
}

const extMan = ".8"
const fileMarkdown = "README.md"

func methodsOf(app *c.App) map[string]method {
	return map[string]method{
		"manpage":  {name: app.Name + extMan, make: app.ToMan},
		"markdown": {name: fileMarkdown, make: app.ToMarkdown},
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

func (m method) run(outdir string) error {
	s, err := m.make()
	if err != nil {
		return err
	}
	fname := filepath.Join(outdir, m.name)
	return ioutil.WriteFile(fname, []byte(s), 0744)
}
