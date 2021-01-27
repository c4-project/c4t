// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"io"

	"github.com/c4-project/c4t/internal/remote"

	"github.com/c4-project/c4t/internal/copier"
	"github.com/c4-project/c4t/internal/plan"
)

// Factory is the interface of factories for machine node runners.
//
// Runner factories can contain disposable state (for example, long-running SSH connections), and so can be closed.
type Factory interface {
	// MakeRunner creates a new Runner, representing a particular invoker session on a machine, and outputting to ldir.
	// It takes the plan p in case the factory is waiting to get machine configuration from it.
	MakeRunner(ldir string, p *plan.Plan, obs ...copier.Observer) (Runner, error)

	// Runner spawners can be closed once no more runners are needed.
	// For SSH runner spawners, this will close the SSH connection.
	io.Closer
}

// FactoryFromRemoteConfig creates a remote factory using gc and mc, if mc is non-nil; else, it creates a local factory.
func FactoryFromRemoteConfig(gc *remote.Config, mc *remote.MachineConfig) (Factory, error) {
	if mc == nil {
		return LocalFactory{}, nil
	}
	return NewRemoteFactory(gc, mc)
}
