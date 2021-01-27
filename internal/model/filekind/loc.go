// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package filekind

// Loc is the bitflag enumeration of file locations.
type Loc uint8

const (
	// InOrig marks that a mapping is part of the original source of a subject.
	InOrig Loc = 1 << iota
	// InFuzz marks that a mapping is part of a fuzz.
	InFuzz
	// InCompile marks that a mapping is part of a compile.
	InCompile
	// InRecipe marks that a mapping is part of a recipe.
	InRecipe
)

// Matches checks whether this location is included in pat.
func (l Loc) Matches(pat Loc) bool {
	return l&pat == l
}
