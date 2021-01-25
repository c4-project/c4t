// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation

import (
	"encoding/csv"
	"sort"
	"strconv"

	"github.com/c4-project/c4t/internal/helper/errhelp"
)

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
			s.Selections[mut]++
			s.Hits[mut] += h.NumHits
			if h.Killed {
				s.Kills[mut]++
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

// Mutants returns a sorted list of all mutant IDs seen in this statset.
func (s *Statset) Mutants() []Mutant {
	muts := make([]Mutant, len(s.Selections))
	i := 0
	for k := range s.Selections {
		// Including mutants that were selected 0 times, because that's interesting.
		muts[i] = k
		i++
	}
	sort.Slice(muts, func(i, j int) bool {
		return muts[i] < muts[j]
	})
	return muts
}

// DumpMutationCSV dumps into w a CSV representation of this mutation statistics set.
// Each line in the record has mid as a prefix.
// The writer is flushed at the end of this dump.
func (s *Statset) DumpCSV(w *csv.Writer, mid string) error {
	var err error
	for _, mut := range s.Mutants() {
		mutstr := strconv.FormatUint(mut, 10)
		selstr := strconv.FormatUint(s.Selections[mut], 10)
		hitstr := strconv.FormatUint(s.Hits[mut], 10)
		killstr := strconv.FormatUint(s.Kills[mut], 10)

		if err = w.Write([]string{mid, mutstr, selstr, hitstr, killstr}); err != nil {
			break
		}
	}

	w.Flush()
	return errhelp.FirstError(err, w.Error())
}
