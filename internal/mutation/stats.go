// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation

import (
	"encoding/csv"
	"sort"
	"strconv"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/helper/errhelp"
)

// Statset holds statistics for each mutant in a mutation testing campaign.
type Statset struct {
	// ByMutant records statsets for each mutant.
	ByMutant map[Mutant]MutantStatset
}

// Reset resets all of this statset's maps to empty, but non-nil.
func (s *Statset) Reset() {
	s.ByMutant = map[Mutant]MutantStatset{}
}

// AddAnalysis adds the information from mutation analysis a to this statset.
func (s *Statset) AddAnalysis(a Analysis) {
	s.ensure()

	for mut, hits := range a {
		m := s.ByMutant[mut]
		m.addAnalysis(hits)
		s.ByMutant[mut] = m
	}
}

func (s *Statset) ensure() {
	if s.ByMutant == nil {
		s.ByMutant = map[Mutant]MutantStatset{}
	}
}

type MutantStatset struct {
	// Selections records the number of times this mutant has been selected.
	Selections uint64
	// Hits records the number of times this mutant has been hit (including kills).
	Hits uint64
	// Kills records the number of selections that resulted in kills.
	Kills uint64
	// Statuses records, for each status, the number of selections that resulted in that status.
	Statuses map[status.Status]uint64
}

func (s *MutantStatset) addAnalysis(hits MutantAnalysis) {
	s.ensure()

	for _, h := range hits {
		s.Selections++
		s.Hits += h.NumHits
		if h.Killed() {
			s.Kills++
		}

		s.Statuses[h.Status]++
	}
}

func (s *MutantStatset) ensure() {
	if s.Statuses == nil {
		s.Statuses = map[status.Status]uint64{}
	}
}

// Mutants returns a sorted list of all mutant IDs seen in this statset.
func (s *Statset) Mutants() []Mutant {
	muts := make([]Mutant, len(s.ByMutant))
	i := 0
	for k := range s.ByMutant {
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
		if err = s.dumpMutant(w, mid, mut, s.ByMutant[mut]); err != nil {
			break
		}
	}

	w.Flush()
	return errhelp.FirstError(err, w.Error())
}

func (s *Statset) dumpMutant(w *csv.Writer, mid string, m Mutant, mstats MutantStatset) error {
	mstats.ensure()
	cells := []string{mid, fint(m), fint(mstats.Selections), fint(mstats.Hits), fint(mstats.Kills)}
	for i := status.Ok; i <= status.Last; i++ {
		cells = append(cells, fint(mstats.Statuses[i]))
	}
	return w.Write(cells)
}

func fint(i uint64) string {
	return strconv.FormatUint(i, 10)
}
