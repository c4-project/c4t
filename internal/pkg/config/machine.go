// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package config

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/remote"
)

// Machine is a config record for a particular machine.
type Machine struct {
	// SSH contains, if present, information about how to dial into a remote machine through SSH.
	SSH *remote.MachineConfig `toml:"ssh,omitempty"`

	model.Machine

	// Compilers contains information about the compilers attached to this machine.
	Compilers map[string]model.Compiler `toml:"compilers,omitempty"`
}
