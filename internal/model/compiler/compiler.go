// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package compiler contains types for compilers, which are a particular type of service.
package compiler

import (
	"github.com/MattWindsor91/act-tester/internal/model/compiler/optlevel"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// Config collects the part of a compiler's specification that comes from the act-tester configuration.
type Config struct {
	// Style is the declared style of the compile.
	Style id.ID `toml:"style"`

	// Arch is the architecture (or 'emits') ID for the compiler.
	Arch id.ID `toml:"arch"`

	// Run contains information on how to run the compiler.
	Run *service.RunInfo `toml:"run,omitempty"`

	// Opt contains information on the optimisation levels to select for the compiler.
	Opt *optlevel.Selection `toml:"opt,omitempty"`
}

// Compiler collects all test-relevant information about a compiler.
type Compiler struct {
	// SelectedOpt refers to an optimisation level chosen using the compiler's configured optimisation selection.
	SelectedOpt *optlevel.Named `toml:"selected_opt,omitempty"`

	Config
}
