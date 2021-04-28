// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package config

import (
	"context"
	"strings"

	"github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/model/service/compiler"

	bimpl "github.com/c4-project/c4t/internal/serviceimpl/backend"
	cimpl "github.com/c4-project/c4t/internal/serviceimpl/compiler"

	"github.com/c4-project/c4t/internal/model/service"

	"github.com/c4-project/c4t/internal/machine"
)

// ProberSet contains the various probers used by configuration probing.
type ProberSet struct {
	// Machine is a machine prober.
	Machine machine.Prober
	// Backend is a backend resolver used for probing.
	Backend backend.Prober
	// Compiler is a compiler prober used for probing.
	Compiler compiler.Prober
}

// LocalProberSet provides a prober set suitable for most local probing.
func LocalProberSet() ProberSet {
	return ProberSet{
		Machine:  machine.LocalProber(),
		Backend:  &bimpl.Resolve,
		Compiler: &cimpl.CResolve,
	}
}

// Probe populates c with information found by scrutinising the current machine.
func (p ProberSet) Probe(ctx context.Context, sr service.Runner, c *Config) error {
	if err := p.probeMachines(ctx, sr, c); err != nil {
		return err
	}
	return p.probeBackends(ctx, sr, c)
}

func (p ProberSet) probeMachines(ctx context.Context, sr service.Runner, c *Config) error {
	if c.RawMachines == nil {
		c.RawMachines = make(map[string]machine.Config)
	}

	hname := hostnameOrDefault(p.Machine)
	var err error
	if _, ok := c.RawMachines[hname]; !ok {
		c.RawMachines[hname], err = p.probeMachine(ctx, sr)
	}

	return err
}

func (p ProberSet) probeBackends(ctx context.Context, sr service.Runner, c *Config) error {
	var err error
	c.Backends, err = p.Backend.Probe(ctx, sr)
	return err
}

func (p ProberSet) probeMachine(ctx context.Context, sr service.Runner) (machine.Config, error) {
	var (
		c   machine.Config
		err error
	)
	if err = c.Probe(p.Machine); err != nil {
		return c, err
	}
	if c.RawCompilers == nil {
		c.RawCompilers, err = p.Compiler.Probe(ctx, sr)
	}

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
