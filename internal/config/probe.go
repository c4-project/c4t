// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package config

import (
	"context"
	"strings"

	"github.com/c4-project/c4t/internal/model/service"
	"github.com/c4-project/c4t/internal/serviceimpl/backend"

	"github.com/c4-project/c4t/internal/machine"
)

// Probe populates this configuration with information found by scrutinising the current machine.
func (c *Config) Probe(ctx context.Context, sr service.Runner, m machine.Prober) error {
	if err := c.probeMachines(m); err != nil {
		return err
	}
	return c.probeBackends(ctx, sr)
}

func (c *Config) probeMachines(m machine.Prober) error {
	if c.RawMachines == nil {
		c.RawMachines = make(map[string]machine.Config)
	}

	hname := hostnameOrDefault(m)
	var err error
	if _, ok := c.RawMachines[hname]; !ok {
		c.RawMachines[hname], err = probeConfig(m)
	}

	return err
}

func (c *Config) probeBackends(ctx context.Context, sr service.Runner) error {
	var err error
	// TODO(@MattWindsor91): should this be hardcoded?
	c.Backends, err = backend.Resolve.Probe(ctx, sr)
	return err
}

func probeConfig(m machine.Prober) (machine.Config, error) {
	var c machine.Config

	err := c.Probe(m)

	return c, err
}

const defaultHostname = "localhost"

func hostnameOrDefault(m machine.Prober) string {
	hname, err := m.Hostname()
	if err != nil {
		return defaultHostname
	}
	hs := strings.Split(hname, ".")
	if len(hs) == 0 {
		return defaultHostname
	}
	return hs[0]
}
