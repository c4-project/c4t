// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package config describes the top-level tester configuration.

package config

import (
	"errors"

	"github.com/MattWindsor91/c4t/internal/model/service/backend"

	"github.com/MattWindsor91/c4t/internal/model/service/fuzzer"

	"github.com/MattWindsor91/c4t/internal/quantity"

	"github.com/MattWindsor91/c4t/internal/machine"

	"github.com/MattWindsor91/c4t/internal/model/id"

	"github.com/MattWindsor91/c4t/internal/remote"
)

// Config is a top-level tester config struct.
type Config struct {
	// Backend contains information about the backend being used to generate recipes.
	Backend *backend.Spec `toml:"backend,omitempty"`

	// Machines enumerates the machines available for testing.
	Machines machine.ConfigMap `toml:"machines,omitempty"`

	// Quantities gives the default quantities for the director.
	Quantities quantity.RootSet `toml:"quantities,omitempty"`

	// Fuzz contains fuzzer config overrides.
	Fuzz *fuzzer.Configuration `toml:"fuzz,omitempty"`

	// SSH contains top-level SSH configuration.
	SSH *remote.Config `toml:"ssh,omitempty"`

	// Paths contains path configuration for the config file.
	Paths Pathset `toml:"paths,omitempty"`
}

// FindBackend uses the configuration to find a backend matching criteria cr.
func (c *Config) FindBackend(cr backend.Criteria) (*backend.NamedSpec, error) {
	// TODO(@MattWindsor91): this needs rearranging a bit.
	if c.Backend == nil {
		return nil, errors.New("backend nil")
	}
	return cr.Find([]backend.NamedSpec{
		{
			ID:   id.ID{},
			Spec: *c.Backend,
		},
	})
}
