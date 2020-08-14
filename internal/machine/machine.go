// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package machine contains models for compiler-bearing machines.
package machine

import (
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/remote"
)

// Machine represents the information about a machine that is relevant to the tester.
type Machine struct {
	// Cores is the number of known cores on the machine.
	// If zero, there is no known core count.
	Cores int `toml:"cores,omitzero" json:"cores,omitempty"`

	// SSH contains, if present, information about how to dial into a remote machine through SSH.
	SSH *remote.MachineConfig `toml:"ssh,omitempty" json:"ssh,omitempty"`
}

// Named wraps a plan machine with its ID.
type Named struct {
	// ID is the ID of the machine.
	ID id.ID `toml:"id,omitzero" json:"id,omitempty"`
	Machine
}
