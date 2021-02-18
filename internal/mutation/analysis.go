// Copyright (c) 2020-2021 C4 Project
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
type Analysis map[Index]MutantAnalysis

// RegisterMutant registers the mutant record m in the analysis.
//
// This is necessary, at the moment, to put things like the mutant's operator
// and variant information in the analysis table.
func (a Analysis) RegisterMutant(m Mutant) {
	ma, ok := a[m.Index]
	if ok {
		return
	}
	ma.Mutant = m
	a[m.Index] = ma
}

// AddCompilation merges any mutant information extracted from log to this analysis.
// Such analysis is filed under compilation name comp, and status determines the status of the compilation.
func (a Analysis) AddCompilation(comp compilation.Name, log string, status status.Status) {
	for i, hits := range ScanLines(strings.NewReader(log)) {
		ma := a[i]
		// Hopefully, we should've pre-registered this mutant's information with the analysis, but this is a failsafe.
		if ma.Mutant.Index == 0 {
			ma.Mutant = AnonMutant(i)
		}
		ma.AddSelection(SelectionAnalysis{
			NumHits: hits,
			Status:  status,
			HitBy:   comp,
		})
		a[i] = ma
	}
}

// Kills determines the mutants that were killed.
func (a Analysis) Kills() []Mutant {
	muts := make([]Mutant, 0, len(a))
	for _, mstat := range a {
		if mstat.Killed {
			muts = append(muts, mstat.Mutant)
		}
	}
	return muts
}

// MutantAnalysis is the type of individual mutant analyses.
type MutantAnalysis struct {
	// Mutant contains the full record for this mutant.
	Mutant Mutant
	// Killed records whether this mutant was killed.
	Killed bool
	// Selections contains the per-selection analysis for this mutant.
	Selections []SelectionAnalysis
}

// AddSelection adds sel to a's selection analyses.
func (a *MutantAnalysis) AddSelection(sel SelectionAnalysis) {
	a.Selections = append(a.Selections, sel)
	a.Killed = a.Killed || sel.Killed()
}

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
