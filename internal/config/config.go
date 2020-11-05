// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package config describes the top-level tester configuration.

package config

import (
	"context"
	"errors"
	"fmt"

	backend2 "github.com/MattWindsor91/act-tester/internal/model/service/backend"

	fuzzer2 "github.com/MattWindsor91/act-tester/internal/model/service/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/machine"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/remote"
)

// Config is a top-level tester config struct.
type Config struct {
	// Backend contains information about the backend being used to generate recipes.
	Backend *backend2.Spec `toml:"backend,omitempty"`

	// Machines enumerates the machines available for testing.
	Machines machine.ConfigMap `toml:"machines,omitempty"`

	// Quantities gives the default quantities for the director.
	Quantities quantity.RootSet `toml:"quantities,omitempty"`

	// Fuzz contains fuzzer config overrides.
	Fuzz *fuzzer2.Configuration `toml:"fuzz,omitempty"`

	// SSH contains top-level SSH configuration.
	SSH *remote.Config `toml:"ssh,omitempty"`

	// Paths contains path configuration for the config file.
	Paths Pathset `toml:"paths,omitempty"`
}

// FindBackend uses the configuration to find a backend with style style.
func (c *Config) FindBackend(_ context.Context, style id.ID, _ ...id.ID) (*backend2.Spec, error) {
	// TODO(@MattWindsor91): this needs rearranging a bit.
	if c.Backend == nil {
		return nil, errors.New("backend nil")
	}
	if !c.Backend.Style.Equal(style) {
		return nil, fmt.Errorf("backend doesn't match given style: got=%q, want=%q", c.Backend.Style.String(), style.String())
	}
	return c.Backend, nil
}
