// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation

import (
	"github.com/c4-project/c4t/internal/subject/compilation"
)

// Analysis is the type of mutation testing analyses.
type Analysis map[uint64]MutantAnalysis

// MutantAnalysis is the type of individual mutant analyses.
type MutantAnalysis []HitAnalysis

// HitAnalysis is the type of analyses for a one or more hits of a mutant by a compilation.
type HitAnalysis struct {
	// NumHits is the number of times this compilation hit the mutant.
	NumHits int `json:"num_hits"`

	// Killed is true provided that this hit resulted in a kill.
	Killed bool `json:"killed"`

	// HitBy is the name of the compilation that hit this mutant.
	HitBy compilation.Name `json:"by"`
}
