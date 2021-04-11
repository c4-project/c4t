// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package machine

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/c4-project/c4t/internal/id"
)

// Prober contains functionality for probing various parts of a machine.
type Prober interface {
	// Hostname gets the full hostname of the machine.
	Hostname() (string, error)

	// NCores gets the number of cores available for the machine.
	NCores() (int, error)

	// Arch gets the architecture ID for the machine.
	Arch() (id.ID, error)
}

//go:generate mockery --name=Prober

// LocalProber gets a prober that probes the local system.
func LocalProber() Prober {
	return localProber{}
}

// localProber implements MachineProber using various local runtime/OS calls.
type localProber struct{}

// Hostname gets the full hostname of the local machine.
func (localProber) Hostname() (string, error) {
	return os.Hostname()
}

// NCores gets the number of cores on the local machine.
func (localProber) NCores() (int, error) {
	return runtime.NumCPU(), nil
}

// Arch gets the architecture using GOARCH.
func (localProber) Arch() (id.ID, error) {
	return ArchOfGOARCH(runtime.GOARCH)
}

// ErrUnsupportedGOARCH occurs when a GOARCH not supported by the prober occurs
var ErrUnsupportedGOARCH = errors.New("local GOARCH not supported by prober")

// ArchOfGOARCH maps from the Go architecture goarch to a C4 architecture id.
// It fails with ErrUnsupportedGOARCH if the architecture isn't supported by C4.
func ArchOfGOARCH(goarch string) (id.ID, error) {
	switch goarch {
	case "386":
		return id.ArchX86, nil
	case "amd64":
		return id.ArchX8664, nil
	case "arm":
		return id.ArchArm, nil
	case "arm64":
		return id.ArchAArch64, nil
	case "ppc64le":
		return id.ArchPPC64LE, nil
	default:
		return id.ID{}, fmt.Errorf("%w: %s", ErrUnsupportedGOARCH, goarch)
	}
}

// Probe populates this config using the prober m.
func (c *Config) Probe(m Prober) error {
	var err error

	if c.Cores, err = m.NCores(); err != nil {
		return err
	}
	if c.Arch, err = m.Arch(); err != nil {
		return err
	}

	return nil
}
