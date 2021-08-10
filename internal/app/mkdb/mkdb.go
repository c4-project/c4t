// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package mkdb implements the c4f-mkdb app.
package mkdb

import (
	"io"

	"github.com/c4-project/c4t/internal/helper/errhelp"

	"github.com/c4-project/c4t/internal/database"

	"github.com/mitchellh/go-wordwrap"

	"github.com/c4-project/c4t/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	// Name is the name of the mkdb binary.
	Name  = "c4t-mkdb"
	usage = "initialises the c4t database"

	readme = `
This program initialises the c4t analysis database.
It should be run before any use of the database, and re-run on a fresh
database when breaking changes to the database sql occur.
`

	FlagDbPath      = "db-file"
	flagDbPathShort = "f"
	usageDbPath     = "path to the SQLite file"
)

// App is the entry point for c4t-config.
func App(outw, errw io.Writer) *c.App {
	a := &c.App{
		Name:        Name,
		Usage:       usage,
		Description: wordwrap.WrapString(readme, 80),
		Flags:       flags(),
		Action:      run,
	}
	return stdflag.SetCommonAppSettings(a, outw, errw)
}

func flags() []c.Flag {
	return []c.Flag{
		&c.PathFlag{
			Name:    FlagDbPath,
			Aliases: []string{flagDbPathShort},
			Usage:   usageDbPath,
			Value:   database.DefaultPath,
		},
	}
}

func run(ctx *c.Context) error {
	db, err := database.Open(ctx.Path(FlagDbPath))
	if err != nil {
		return err
	}

	derr := database.Create(ctx.Context, db)
	cerr := db.Close()
	return errhelp.FirstError(derr, cerr)
}
