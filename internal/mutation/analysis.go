// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation

import (
	"strings"

	"github.com/c4-project/c4t/internal/subject/compilation"
)

// Analysis is the type of mutation testing analyses.
type Analysis map[Mutant]MutantAnalysis

// AddCompilation merges any mutant information extracted from log to this analysis.
// Such analysis is filed under compilation name comp, and killer determines whether hit mutants was killed.
func (a Analysis) AddCompilation(comp compilation.Name, log string, killer bool) {
	for mut, hits := range ScanLines(strings.NewReader(log)) {
		ana := HitAnalysis{
			// TODO(@MattWindsor91): get rid of this cast somehow
			NumHits: hits,
			Killed:  killer && hits != 0,
			HitBy:   comp,
		}
		a[mut] = append(a[mut], ana)
	}
}

// MutantAnalysis is the type of individual mutant analyses.
type MutantAnalysis []HitAnalysis

// HitAnalysis is the type of analyses for a one or more hits of a mutant by a compilation.
type HitAnalysis struct {
	// NumHits is the number of times this compilation hit the mutant.
	// If this is 0, the mutant was selected but never hit.
	NumHits uint64 `json:"num_hits"`

	// Killed is true provided that this hit resulted in a kill.
	// If the compilation failed, this will be true unless the mutant was never hit (NumHits == 0).
	Killed bool `json:"killed"`

	// HitBy is the name of the compilation that hit this mutant.
	HitBy compilation.Name `json:"by"`
}
