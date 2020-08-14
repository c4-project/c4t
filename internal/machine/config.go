// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package machine

import (
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/MattWindsor91/act-tester/internal/quantity"
)

// Config is a config record for a particular machine.
type Config struct {
	Machine

	// Compilers contains information about the compilers attached to this machine.
	Compilers map[string]compiler.Compiler `toml:"compilers,omitempty"`

	// Quantities contains quantity overrides for this machine.
	Quantities quantity.MachineSet `toml:"quantities,omitempty,omitzero"`
}
