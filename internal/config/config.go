// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package config describes the top-level tester configuration.
package config

import (
	"io"

	"github.com/c4-project/c4t/internal/model/service/backend"
	"github.com/pelletier/go-toml"

	"github.com/c4-project/c4t/internal/model/service/fuzzer"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/machine"

	"github.com/c4-project/c4t/internal/remote"
)

// Config is a top-level tester config struct.
type Config struct {
	// Backends contains information about the backends available for generating recipes.
	//
	// These are given as a list rather than a map because ordering matters: the director will satisfy requests for
	// backends by trying each backend specification in order.
	Backends []backend.NamedSpec `toml:"backends,omitempty"`

	// Machines enumerates the machines available for testing.
	Machines machine.ConfigMap `toml:"machines,omitempty"`

	// Quantities gives the default quantities for the director.
	Quantities quantity.RootSet `toml:"quantities,omitempty"`

	// Fuzz contains fuzzer config overrides.
	Fuzz *fuzzer.Config `toml:"fuzz,omitempty"`

	// SSH contains top-level SSH configuration.
	SSH *remote.Config `toml:"ssh,omitempty"`

	// Paths contains path configuration for the config file.
	Paths Pathset `toml:"paths,omitempty"`
}

// FindBackend uses the configuration to find a backend matching criteria cr.
func (c *Config) FindBackend(cr backend.Criteria) (*backend.NamedSpec, error) {
	return cr.Find(c.Backends)
}

// Dump dumps this configuration to the writer w.
func (c *Config) Dump(w io.Writer) error {
	return toml.NewEncoder(w).Encode(c)
}

// OverrideQuantities is shorthand for overriding the quantity set in this config.
func (c *Config) OverrideQuantities(qs quantity.RootSet) {
	c.Quantities.Override(qs)
}

// DisableFuzz is shorthand for setting this config's fuzzer disabled flag to false.
func (c *Config) DisableFuzz() {
	if c.Fuzz == nil {
		c.Fuzz = &fuzzer.Config{}
	}
	c.Fuzz.Disabled = true
}

// OverrideInputs is shorthand for setting this config's inputs to files, if non-empty.
func (c *Config) OverrideInputs(files []string) error {
	// TODO(@MattWindsor91): push this into pathset?
	files, err := c.Paths.FallbackToInputs(files)
	if err != nil {
		return err
	}
	c.Paths.Inputs = files
	return nil
}
