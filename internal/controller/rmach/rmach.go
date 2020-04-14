// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package rmach handles invoking the 'mach' controller in a network-transparent manner with the act-tester-mach binary.
package rmach

import (
	"github.com/MattWindsor91/act-tester/internal/model/plan"
	"github.com/MattWindsor91/act-tester/internal/remote"
)

// RMach runs the machine-runner, through SSH if needed.
type RMach struct {
	conf   *Config
	plan   *plan.Plan
	runner Runner
}

// New constructs a new Mach with ssh configuration ssh (if any) and local directory dir.
func New(c *Config, p *plan.Plan) (*RMach, error) {
	if err := check(c, p); err != nil {
		return nil, err
	}

	r, err := newRunner(c.Observers.Copy, c.DirLocal, c.SSH, p.Machine.SSH)
	if err != nil {
		return nil, err
	}
	return &RMach{conf: c, plan: p, runner: r}, nil
}

func check(c *Config, p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	return checkConfig(c)
}

func checkConfig(c *Config) error {
	if c == nil {
		return ErrConfigNil
	}
	return c.Check()
}

func newRunner(o []remote.CopyObserver, dir string, c *remote.Config, ssh *remote.MachineConfig) (Runner, error) {
	if ssh == nil {
		return NewLocalRunner(dir), nil
	}
	sc, err := ssh.MachineRunner(c)
	if err != nil {
		return nil, err
	}
	return NewSSHRunner(sc, o, dir), nil
}
