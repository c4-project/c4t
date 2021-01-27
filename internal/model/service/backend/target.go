// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend

// Target is the enumeration of targets of a backend.
type Target uint8

const (
	// ToDefault states that the lifting backend should perform its default lifting.
	ToDefault Target = iota
	// ToStandalone states that the backend should run in a standalone manner in the running phase, without compilation.
	// The backend can still produce files and include them in the recipe, but should not produce instructions.
	ToStandalone
	// ToObjRecipe states that the backend should produce a recipe that emits one or more object files.
	ToObjRecipe
	// ToExeRecipe states that the backend should produce a recipe that emits a compilable executable.
	ToExeRecipe
)

//go:generate stringer -type Target
