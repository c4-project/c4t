// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stat

import (
	"encoding/csv"
	"sort"
	"strconv"

	"github.com/c4-project/c4t/internal/mutation"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/helper/errhelp"
)

// Mutation holds statistics for each mutant in a mutation testing campaign.
type Mutation struct {
	// ByMutant records statsets for each mutant.
	ByMutant map[mutation.Mutant]Mutant
}

// Reset resets all of this statset's maps to empty, but non-nil.
func (m *Mutation) Reset() {
	m.ByMutant = map[mutation.Mutant]Mutant{}
}

// AddAnalysis adds the information from mutation analysis a to this statset.
func (m *Mutation) AddAnalysis(a mutation.Analysis) {
	m.ensure()

	for mut, hits := range a {
		ms := m.ByMutant[mut]
		ms.addAnalysis(hits)
		m.ByMutant[mut] = ms
	}
}

func (m *Mutation) ensure() {
	if m.ByMutant == nil {
		m.ByMutant = map[mutation.Mutant]Mutant{}
	}
}

// Mutants returns a sorted list of all mutant IDs seen in this statset.
func (m *Mutation) Mutants() []mutation.Mutant {
	return m.MutantsWhere(FilterAllMutants)
}

// KilledMutants returns a sorted list of all mutant IDs killed in this statset.
func (m *Mutation) KilledMutants() []mutation.Mutant {
	return m.MutantsWhere(FilterKilledMutants)
}

// MutantFilter is the type of mutant filtering predicates.
type MutantFilter func(m Mutant) bool

var (
	// FilterAllMutants is a mutant filter that allows all mutants.
	FilterAllMutants MutantFilter = func(mutant Mutant) bool { return true }
	// FilterHitMutants is a mutant filter that allows hit mutants only.
	FilterHitMutants MutantFilter = func(mutant Mutant) bool { return 0 < mutant.Hits }
	// FilterKilledMutants is a mutant filter that allows killed mutants only.
	FilterKilledMutants MutantFilter = func(mutant Mutant) bool { return 0 < mutant.Kills }
	// FilterEscapedMutants is a mutant filter that allows only mutants that were hit but not killed.
	FilterEscapedMutants MutantFilter = func(mutant Mutant) bool {
		return 0 < mutant.Hits && 0 == mutant.Kills
	}
)

// MutantsWhere returns a sorted list of mutants satisfying pred.
// (It is a value receiver method to allow calling through templates.)
func (m Mutation) MutantsWhere(pred func(m Mutant) bool) []mutation.Mutant {
	muts := make([]mutation.Mutant, 0, len(m.ByMutant))
	for k, mstat := range m.ByMutant {
		if pred(mstat) {
			muts = append(muts, k)
		}
	}
	sort.Slice(muts, func(i, j int) bool {
		return muts[i] < muts[j]
	})
	return muts
}

// DumpMutationCSV dumps into w a CSV representation of this mutation statistics set.
// Each line in the record has mid as a prefix.
// The writer is flushed at the end of this dump.
func (m *Mutation) DumpCSV(w *csv.Writer, mid string) error {
	var err error
	for _, mut := range m.Mutants() {
		if err = m.dumpMutant(w, mid, mut, m.ByMutant[mut]); err != nil {
			break
		}
	}

	w.Flush()
	return errhelp.FirstError(err, w.Error())
}

// Mutant gives statistics for a particular mutant.
type Mutant struct {
	// Selections records the number of times this mutant has been selected.
	Selections uint64
	// Hits records the number of times this mutant has been hit (including kills).
	Hits uint64
	// Kills records the number of selections that resulted in kills.
	Kills uint64
	// Statuses records, for each status, the number of selections that resulted in that status.
	Statuses map[status.Status]uint64
}

func (m *Mutant) addAnalysis(hits mutation.MutantAnalysis) {
	m.ensure()

	for _, h := range hits {
		m.Selections++
		m.Hits += h.NumHits
		if h.Killed() {
			m.Kills++
		}

		m.Statuses[h.Status]++
	}
}

func (m *Mutant) ensure() {
	if m.Statuses == nil {
		m.Statuses = map[status.Status]uint64{}
	}
}

func (m *Mutation) dumpMutant(w *csv.Writer, machname string, mut mutation.Mutant, mstats Mutant) error {
	mstats.ensure()
	cells := []string{machname, fint(mut), fint(mstats.Selections), fint(mstats.Hits), fint(mstats.Kills)}
	for i := status.Ok; i <= status.Last; i++ {
		cells = append(cells, fint(mstats.Statuses[i]))
	}
	return w.Write(cells)
}

func fint(i uint64) string {
	return strconv.FormatUint(i, 10)
}
