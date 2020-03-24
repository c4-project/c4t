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

	"github.com/MattWindsor91/act-tester/internal/pkg/model/service"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/transfer/remote"
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
	Backend *service.Backend `toml:"backend,omitempty"`

	// Machines enumerates the machines available for testing.
	Machines map[string]Machine `toml:"machines,omitempty"`

	// Quantities gives the default quantities for the director.
	Quantities QuantitySet `toml:"quantities,omitempty"`

	// SSH contains top-level SSH configuration.
	SSH *remote.Config `toml:"ssh,omitempty"`

	// OutDir is the output directory for fully directed test runs.
	OutDir string `toml:"out_dir"`
}

// FindBackend uses the configuration to find a backend with style style.
func (c *Config) FindBackend(_ context.Context, style id.ID, _ ...id.ID) (*service.Backend, error) {
	// TODO(@MattWindsor91): this needs rearranging a bit.
	if c.Backend == nil {
		return nil, errors.New("backend nil")
	}
	if !c.Backend.Style.Equal(style) {
		return nil, fmt.Errorf("backend doesn't match given style: got=%q, want=%q", c.Backend.Style.String(), style.String())
	}
	return c.Backend, nil
}

// MachineIDs gets a sorted slice of machine IDs present in the config.
// It returns an error if any of the configured machines have an invalid ID.
func (c *Config) MachineIDs() ([]id.ID, error) {
	return id.MapKeys(c.Machines)
}

// ListCompilers implements the compiler listing operation using a config.
func (c *Config) ListCompilers(_ context.Context, mid id.ID) (map[string]service.Compiler, error) {
	mstr := mid.String()
	m, ok := c.Machines[mstr]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoMachine, mstr)
	}
	return m.Compilers, nil
}
