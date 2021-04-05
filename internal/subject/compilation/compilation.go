// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compilation

import "github.com/c4-project/c4t/internal/id"

// Compilation holds information about one testcase-compiler pairing.
type Compilation struct {
	// Compile contains any information about the compile phase of this compilation.
	Compile *CompileResult `json:"compile,omitempty"`
	// Run contains any information about the run phase of this compilation.
	Run *RunResult `json:"run,omitempty"`
}

// Map is shorthand for a map from compiler IDs to compilations.
type Map map[id.ID]Compilation
