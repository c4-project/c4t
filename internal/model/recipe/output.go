// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package recipe

// Output is the enumeration of types of recipe outputs.
type Output uint8

const (
	// This recipe doesn't output anything, and exists solely to feed into another recipe.
	OutNothing Output = iota
	// This recipe outputs a file that should be read into the observation parser.
	OutText
	// This recipe outputs an executable that should be run with its output piped into the observation parser.
	OutExe
)

//go:generate stringer -type Output
