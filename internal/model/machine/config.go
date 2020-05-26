// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package machine

import (
	"github.com/MattWindsor91/act-tester/internal/model/compiler"
)

// Config is a config record for a particular machine.
type Config struct {
	Machine

	// Compilers contains information about the compilers attached to this machine.
	Compilers map[string]compiler.Config `toml:"compilers,omitempty"`
}
