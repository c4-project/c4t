// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package mach provides director interoperability with the act-tester-mach binary.
package mach

import (
	"github.com/MattWindsor91/act-tester/internal/director/observer"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/transfer/remote"
)

// Mach runs the machine-runner, through SSH if needed.
type Mach struct {
	// observers is the set of observers to which we are sending updates from the machine-runner.
	observers []builder.Observer

	// runner describes how to run the machine-runner binary.
	runner Runner
}

// New constructs a new Mach with ssh configuration ssh (if any) and local directory dir.
func New(obs []observer.Instance, dir string, c *remote.Config, ssh *remote.MachineConfig) (*Mach, error) {
	r, err := newRunner(observer.LowerToCopy(obs), dir, c, ssh)
	if err != nil {
		return nil, err
	}
	return &Mach{observers: observer.LowerToBuilder(obs), runner: r}, nil
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