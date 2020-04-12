// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package config

import (
	"github.com/MattWindsor91/act-tester/internal/model/compiler"
	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Machine is a config record for a particular machine.
type Machine struct {
	plan.Machine

	// Compilers contains information about the compilers attached to this machine.
	Compilers map[string]compiler.Config `toml:"compilers,omitempty"`
}
