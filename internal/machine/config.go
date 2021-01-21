// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package machine

import (
	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/mutation"
)

// Config is a config record for a particular machine.
//
// The difference between a Machine and a Config is that the latter contains raw configuration data for things that get
// mapped into expanded forms in the plan, for instance compilers.
type Config struct {
	Machine

	// Compilers contains information about the compilers attached to this machine.
	Compilers map[string]compiler.Compiler `toml:"compilers,omitempty"`

	// Mutation contains information about how to mutation-test on this machine.
	Mutation *mutation.Config `toml:"mutation,omitempty"`
}
