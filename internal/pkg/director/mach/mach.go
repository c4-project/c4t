// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package mach provides director interoperability with the act-tester-mach binary.
package mach

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/pkg/transfer/remote"
)

// Mach runs the machine-runner, through SSH if needed.
type Mach struct {
	// observer is the observer to which we are sending updates from the machine-runner.
	observer builder.Observer

	// runner describes how to run the machine-runner binary.
	runner Runner
}

// New constructs a new Mach with ssh configuration ssh (if any) and local directory dir.
func New(o builder.Observer, dir string, c *remote.Config, ssh *remote.MachineConfig) (*Mach, error) {
	m := Mach{observer: o}
	if ssh == nil {
		m.runner = NewLocalRunner(dir)
		return &m, nil
	}

	sc, err := ssh.MachineRunner(c)
	if err != nil {
		return nil, err
	}
	m.runner = NewSSHRunner(sc)
	return &m, nil
}
