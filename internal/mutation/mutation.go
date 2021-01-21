// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package mutation contains support for mutation testing using c4t.
package mutation

// TODO(@MattWindsor91): move these
const (
	// MutantHitPrefix is the prefix of lines from compilers specifying that a mutant has been hit.
	MutantHitPrefix = "MUTATION HIT:"
	// MutantSelectPrefix is the prefix of lines from compilers specifying that a mutant has been selected.
	MutantSelectPrefix = "MUTATION SELECTED:"
)
