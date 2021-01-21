// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation

// Statset holds statistics for each mutant in a mutation testing campaign.
type Statset struct {
	// Selections records the number of times each mutant has been selected.
	Selections map[Mutant]uint64
	// Hits records the number of times each mutant has been hit (including kills).
	Hits map[Mutant]uint64
	// Hits records the number of times each mutant has been hit (including kills).
	Kills map[Mutant]uint64
}

// Reset resets all of this statset's maps to empty, but non-nil.
func (s *Statset) Reset() {
	s.Selections = map[Mutant]uint64{}
	s.Hits = map[Mutant]uint64{}
	s.Kills = map[Mutant]uint64{}
}

// AddAnalysis adds the information from mutation analysis a to this statset.
func (s *Statset) AddAnalysis(a Analysis) {
	s.ensure()

	for mut, hits := range a {
		for _, h := range hits {
			s.Selections[mut] += 1
			s.Hits[mut] += h.NumHits
			if h.Killed {
				s.Kills[mut] += h.NumHits
			}
		}
	}
}

func (s *Statset) ensure() {
	if s.Selections == nil {
		s.Selections = map[Mutant]uint64{}
	}
	if s.Hits == nil {
		s.Hits = map[Mutant]uint64{}
	}
	if s.Kills == nil {
		s.Kills = map[Mutant]uint64{}
	}
}
