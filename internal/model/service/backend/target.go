// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package backend

// Target is the enumeration of targets of a backend.
type Target uint8

const (
	// ToDefault states that the lifting backend should return whatever it feels most comfortable with.
	ToDefault Target = iota
	// ToStandalone states that the backend should run to completion without further compilation or running;
	// it will emit a single text output file and stub recipe.
	ToStandalone
	// ToObjRecipe states that the backend should produce a recipe that emits one or more object files.
	ToObjRecipe
	// ToExeRecipe states that te backend should produce a recipe that emits a compilable executable.
	ToExeRecipe
)

//go:generate stringer -type Target
