// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compilation

// Compilation holds information about one testcase-compiler pairing.
type Compilation struct {
	// Compile contains any information about the compile phase of this compilation.
	Compile *CompileResult `json:"compile,omitempty"`
	// Run contains any information about the run phase of this compilation.
	Run *RunResult `json:"run,omitempty"`
}
