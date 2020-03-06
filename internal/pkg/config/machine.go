// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package config

import "github.com/MattWindsor91/act-tester/internal/pkg/model"

// Machine is a config record for a particular machine.
type Machine struct {
	// SSH contains, if present, information about how to dial into a remote machine through SSH.
	SSH *MachineSSH `toml:"ssh,omitempty"`

	model.Machine

	// Compilers contains information about the compilers attached to this machine.
	Compilers map[string]model.Compiler `toml:"compilers,omitempty"`
}

// MachineSSH is SSH configuration for a remote machine.
type MachineSSH struct {
	// The host to use when dialing into the machine.
	Host string `toml:"host"`
	// The user to use when dialing into the machine.
	User string `toml:"user,omitzero"`
	// The directory to which we shall copy intermediate files.
	DirCopy string `toml:"copy_dir"`
}
