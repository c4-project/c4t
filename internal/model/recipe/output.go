// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package recipe

// Output is the enumeration of types of recipe outputs.
type Output uint8

const (
	// This recipe doesn't output anything, and so need not be run.
	// The recipe can still include files to make available to the backend's run-time driver.
	OutNothing Output = iota
	// This recipe outputs an object file that should be fed into another recipe.
	OutObj
	// This recipe outputs an executable that should be run with its output piped into the observation parser.
	OutExe
)

//go:generate stringer -type Output
