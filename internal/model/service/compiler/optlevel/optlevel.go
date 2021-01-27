// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package optlevel contains types that capture information about compiler optimisation levels.
package optlevel

// Level holds information about an optimisation level.
type Level struct {
	// Optimises is true if this optimisation level actually optimises in any sense.
	// Counter-examples include '-O0' in gcc, and '/Od' in MSVC.
	Optimises bool `toml:"optimises,omitempty" json:"optimises,omitempty"`

	// Bias is this optimisation level's bias, if known.
	Bias Bias `toml:"bias,omitempty" json:"bias,omitempty"`

	// BreaksStandards is true if this optimisation level plays hard and fast with the C standard.
	// Examples include '-Ofast' in GCC.
	BreaksStandards bool `toml:"breaks_standards" json:"breaks_standards"`
}

// Named wraps a Level with its command-line name.
type Named struct {
	// Name is the name of the optimisation level, which should be its command-line invocation less any prefix.
	Name string `toml:"name" json:"name"`

	Level
}
