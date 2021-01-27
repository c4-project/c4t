// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package machine

import (
	"os"
	"runtime"
)

// Prober contains functionality for probing various parts of a machine.
type Prober interface {
	// Hostname gets the full hostname of the machine.
	Hostname() (string, error)

	// NCores gets the number of cores available for the machine.
	NCores() (int, error)
}

//go:generate mockery --name=Prober

// LocalProber implements MachineProber using various local runtime/OS calls.
type LocalProber struct{}

// Hostname gets the full hostname of the local machine.
func (l LocalProber) Hostname() (string, error) {
	return os.Hostname()
}

// NCores gets the number of cores on the local machine.
func (l LocalProber) NCores() (int, error) {
	return runtime.NumCPU(), nil
}

// Probe populates the machine c using the prober m.
func Probe(c *Machine, m Prober) error {
	var err error

	if c.Cores, err = m.NCores(); err != nil {
		return err
	}

	return nil
}
