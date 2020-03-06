// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package director contains the top-level ACT test director, which manages a full testing campaign.
package director

import (
	"context"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
)

// Director contains the main state and configuration for the test director.
type Director struct {
	// config is the configuration for the director.
	config *Config

	// l is the logger for the director.
	l *log.Logger

	// paths tells the director where to store its
	paths *Pathset
}

// New creates a new Director given a global act-tester config.
// It fails if the config is missing or ill-formed.
func New(c *Config) (*Director, error) {
	if err := checkConfig(c); err != nil {
		return nil, err
	}

	return &Director{config: c, l: iohelp.EnsureLog(c.Logger)}, nil
}

func checkConfig(c *Config) error {
	if c == nil {
		return ErrConfigNil
	}
	if c.Paths == nil {
		return iohelp.ErrPathsetNil
	}
	return nil
}

// Direct runs the director d.
func (d *Director) Direct(_ context.Context) error {
	d.l.Print("making directories")
	if err := d.paths.Prepare(); err != nil {
		return err
	}

	// TODO(@MattWindsor91)
	return nil
}
