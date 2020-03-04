// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package model

// Compiler collects the test-relevant information about a compiler.
type Compiler struct {
	// Style is the declared style of the backend.
	Style ID `toml:"style"`

	// Arch is the architecture (or 'emits') CompilerID for the compiler.
	Arch ID `toml:"arch"`
}

// NamedCompiler wraps a Compiler with its ID.
type NamedCompiler struct {
	// ID is the ID of the compiler.
	ID ID `toml:"id"`

	Compiler
}
