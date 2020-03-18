// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package model

import "github.com/MattWindsor91/act-tester/internal/pkg/model/id"

// Compiler collects the test-relevant information about a compiler.
type Compiler struct {
	// Style is the declared style of the backend.
	Style id.ID `toml:"style"`

	// Arch is the architecture (or 'emits') ID for the compiler.
	Arch id.ID `toml:"arch"`

	// Run contains information on how to run the compiler.
	Run *CompilerRunInfo `toml:"run,omitempty"`
}

// NamedCompiler wraps a Compiler with its ID.
type NamedCompiler struct {
	// ID is the ID of the compiler.
	ID id.ID `toml:"id"`

	Compiler
}

// CompilerRunInfo gives hints as to how to run a compiler.
type CompilerRunInfo struct {
	// Cmd overrides the command for the compiler.
	Cmd string `toml:"cmd,omitzero"`

	// Args specifies (extra) arguments to supply to the compiler.
	Args []string `toml:"args,omitempty"`
}

// Override creates run information by overlaying this run information with that in new.
func (c CompilerRunInfo) Override(new *CompilerRunInfo) CompilerRunInfo {
	if new == nil {
		return c
	}
	return CompilerRunInfo{
		Cmd:  overrideCmd(c.Cmd, new.Cmd),
		Args: append(c.Args, new.Args...),
	}
}

func overrideCmd(old, new string) string {
	if new == "" {
		return old
	}
	return new
}
