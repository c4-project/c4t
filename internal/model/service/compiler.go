// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package service

import "github.com/MattWindsor91/act-tester/internal/model/id"

// Compiler collects the test-relevant information about a compiler.
type Compiler struct {
	// Style is the declared style of the compile.
	Style id.ID `toml:"style"`

	// Arch is the architecture (or 'emits') ID for the compiler.
	Arch id.ID `toml:"arch"`

	// Run contains information on how to run the compiler.
	Run *RunInfo `toml:"run,omitempty"`
}

// NamedCompiler wraps a Compiler with its ID.
type NamedCompiler struct {
	// ID is the ID of the compiler.
	ID id.ID `toml:"id"`

	Compiler
}
