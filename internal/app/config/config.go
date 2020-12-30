// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package config implements the act-config app.
package config

import (
	"io"

	"github.com/c4-project/c4t/internal/config"
	"github.com/c4-project/c4t/internal/machine"
	"github.com/c4-project/c4t/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	// Name is the name of the config binary.
	Name  = "c4t-config"
	usage = "initialises config"

	readme = `
   This program produces an initial c4 config file for the current system.
`
)

// App is the entry point for c4t-config.
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
	return stdflag.SetCommonAppSettings(a, outw, errw)
}

func flags() []c.Flag {
	return []c.Flag{}
}

func run(ctx *c.Context, outw io.Writer, _ io.Writer) error {
	cfg := config.Config{}
	if err := cfg.Probe(machine.LocalProber{}); err != nil {
		return err
	}
	return cfg.Dump(outw)
}
