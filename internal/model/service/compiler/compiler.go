// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package compiler contains types for compilers, which are a particular type of service.
package compiler

import (
	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/model/service"
	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"
)

// Compiler represents fully prepared baseline configuration for a compiler.
//
// The distinction between Compiler and Config is that the latter represents a raw form of a Compiler coming from a
// config file; the distinction between Compiler and Instance is that the latter is the former plus a set of decisions
// on the specific invocation the compiler will have for a test run (eg, optimisation levels, mutant indices, etc).
type Compiler struct {
	// Disabled specifies whether this compiler has been disabled.
	Disabled bool `toml:"disabled,omitempty" json:"disabled,omitempty"`

	// Style is the declared style of the compile.
	Style id.ID `toml:"style" json:"style"`

	// Arch is the architecture (or 'emits') ID for the compiler.
	Arch id.ID `toml:"arch" json:"arch"`

	// Run contains information on how to run the compiler.
	Run *service.RunInfo `toml:"run,omitempty" json:"run,omitempty"`

	// MOpt contains information on the 'mopt' (compiler architecture tuning) levels to select for the compiler.
	MOpt *optlevel.Selection `toml:"mopt,optempty" json:"mopt,omitempty"`

	// Opt contains information on the optimisation levels to select for the compiler.
	Opt *optlevel.Selection `toml:"opt,omitempty" json:"opt,omitempty"`
}

// Config denotes raw configuration for a Compiler.
type Config Compiler
