// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stat

import (
	"github.com/c4-project/c4t/internal/director"
	"github.com/c4-project/c4t/internal/mutation"
	"github.com/c4-project/c4t/internal/plan/analysis"
	"github.com/c4-project/c4t/internal/subject/status"
)

// Machine is a statistics set for a specific machine.
type Machine struct {
	// TotalCycles is the number of the cycles the machine has accrued.
	TotalCycles uint64 `json:"total_cycles"`
	// LastCycle is the last announced cycle in this session.
	LastCycle director.Cycle `json:"last_cycle,omitempty"`

	// StatusTotals contains status totals since this statset started.
	StatusTotals map[status.Status]uint64 `json:"status_totals"`
	// SessionStatusTotals contains status totals since this session started.
	// It may be empty if this machine has not yet been active this session.
	SessionStatusTotals map[status.Status]uint64 `json:"session_status_totals,omitempty"`

	// Mutation contains totals for mutation testing since this statset started.
	Mutation mutation.Statset `json:"mutation,omitempty"`
	// Mutation contains totals for mutation testing since this session started.
	SessionMutation mutation.Statset `json:"session_mutation,omitempty"`
}

// ResetForSession removes from this statset any statistics that no longer apply across session boundaries.
func (m *Machine) ResetForSession() {
	m.LastCycle = director.Cycle{}
	m.SessionStatusTotals = make(map[status.Status]uint64)
	m.SessionMutation.Reset()
}

// AddAnalysis adds the information from analysis a to this machine statset.
func (m *Machine) AddAnalysis(a analysis.Analysis) {
	m.addStatusTotals(a)
	m.addMutation(a)
}

func (m *Machine) addStatusTotals(a analysis.Analysis) {
	if m.StatusTotals == nil {
		m.StatusTotals = make(map[status.Status]uint64)
	}
	if m.SessionStatusTotals == nil {
		m.SessionStatusTotals = make(map[status.Status]uint64)
	}
	for i := status.Ok; i <= status.Last; i++ {
		count := uint64(len(a.ByStatus[i]))
		m.StatusTotals[i] += count
		m.SessionStatusTotals[i] += count
	}
}

func (m *Machine) addMutation(a analysis.Analysis) {
	if len(a.Mutation) == 0 {
		return
	}
	m.Mutation.AddAnalysis(a.Mutation)
	m.SessionMutation.AddAnalysis(a.Mutation)
}
