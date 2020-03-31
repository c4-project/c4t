// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan

// Machine represents the information about a machine that is relevant to the tester.
type Machine struct {
	// Cores is the number of known cores on the machine.
	// If zero, there is no known core count.
	Cores int `toml:"cores,omitzero"`
}
