// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package machine

import (
	"github.com/c4-project/c4t/internal/model/service/compiler"
)

// Config is a config record for a particular machine.
type Config struct {
	Machine

	// Compilers contains information about the compilers attached to this machine.
	Compilers map[string]compiler.Compiler `toml:"compilers,omitempty"`
}
