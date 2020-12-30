// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package config

import (
	"strings"

	"github.com/c4-project/c4t/internal/machine"
)

// Probe populates this configuration with information found by scrutinising the current machine.
func (c *Config) Probe(m machine.Prober) error {
	if err := c.probeMachines(m); err != nil {
		return err
	}
	return nil
}

func (c *Config) probeMachines(m machine.Prober) error {
	if c.Machines == nil {
		c.Machines = make(machine.ConfigMap)
	}

	hname := hostnameOrDefault(m)
	var err error
	if _, ok := c.Machines[hname]; !ok {
		c.Machines[hname], err = probeConfig(m)
	}

	return err
}

func probeConfig(m machine.Prober) (machine.Config, error) {
	var (
		c   machine.Config
		err error
	)

	err = machine.Probe(&(c.Machine), m)

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
