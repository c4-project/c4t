// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package config

import (
	"github.com/MattWindsor91/act-tester/internal/model/plan"
	"github.com/MattWindsor91/act-tester/internal/model/service"
	"github.com/MattWindsor91/act-tester/internal/transfer/remote"
)

// Machine is a config record for a particular machine.
type Machine struct {
	// SSH contains, if present, information about how to dial into a remote machine through SSH.
	SSH *remote.MachineConfig `toml:"ssh,omitempty"`

	plan.Machine

	// Compilers contains information about the compilers attached to this machine.
	Compilers map[string]service.Compiler `toml:"compilers,omitempty"`
}
