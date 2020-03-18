// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package config describes the top-level tester configuration.

// TODO(@MattWindsor91): slowly wrest control of the configuration from OCaml act.

package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/remote"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

var (
	// ErrNoMachine occurs when we try to look up the compilers of a missing machine.
	ErrNoMachine = errors.New("no such machine")

	// ErrNil occurs when we try to build something using a nil config.
	ErrNil = errors.New("config nil")
)

// Config is a top-level tester config struct.
type Config struct {
	// Backend contains information about the backend being used to generate test harnesses.
	Backend *model.Backend `toml:"backend,omitempty"`

	// Machines enumerates the machines available for testing.
	Machines map[string]Machine `toml:"machines,omitempty"`

	// Quantities gives the default quantities for the director.
	Quantities QuantitySet `toml:"quantities,omitempty"`

	// SSH contains top-level SSH configuration.
	SSH *remote.Config `toml:"ssh,omitempty"`

	// OutDir is the output directory for fully directed test runs.
	OutDir string `toml:"out_dir"`
}

// MachineIDs gets a sorted slice of machine IDs present in the config.
// It returns an error if any of the configured machines have an invalid ID.
func (c *Config) MachineIDs() ([]id.ID, error) {
	var (
		err error
		i   int
	)

	mids := make([]id.ID, len(c.Machines))
	for mstr := range c.Machines {
		if mids[i], err = id.TryFromString(mstr); err != nil {
			return nil, err
		}
		i++
	}
	return mids, nil
}

// ListCompilers implements the compiler listing operation using a config.
func (c *Config) ListCompilers(_ context.Context, mid id.ID) (map[string]model.Compiler, error) {
	mstr := mid.String()
	m, ok := c.Machines[mstr]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoMachine, mstr)
	}
	return m.Compilers, nil
}
