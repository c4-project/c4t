// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation

import (
	"strings"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/subject/compilation"
)

// Analysis is the type of mutation testing analyses.
type Analysis map[Mutant]MutantAnalysis

// AddCompilation merges any mutant information extracted from log to this analysis.
// Such analysis is filed under compilation name comp, and status determines the status of the compilation.
func (a Analysis) AddCompilation(comp compilation.Name, log string, status status.Status) {
	for mut, hits := range ScanLines(strings.NewReader(log)) {
		a[mut] = append(a[mut], SelectionAnalysis{
			NumHits: hits,
			Status:  status,
			HitBy:   comp,
		})
	}
}

// HasKills determines whether there is at least one killed mutant in this analysis.
func (a Analysis) HasKills() bool {
	for _, mut := range a {
		for _, hit := range mut {
			if hit.Killed() {
				return true
			}
		}
	}
	return false
}

// MutantAnalysis is the type of individual mutant analyses.
type MutantAnalysis []SelectionAnalysis

// SelectionAnalysis represents one instance where a compilation selected a particular mutant.
type SelectionAnalysis struct {
	// NumHits is the number of times this compilation hit the mutant.
	// If this is 0, the mutant was selected but never hit.
	NumHits uint64 `json:"num_hits"`

	// Status was the main status of the compilation, which determines whether the selection killed the mutant.
	Status status.Status `json:"status"`

	// HitBy is the name of the compilation that hit this mutant.
	HitBy compilation.Name `json:"by"`
}

// Killed gets whether this selection resulted in a kill (hit at least once and resulted in a flagged status).
func (h SelectionAnalysis) Killed() bool {
	return 0 < h.NumHits && h.Status == status.Flagged
}
