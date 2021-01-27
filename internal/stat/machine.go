// Copyright (c) 2020-2021 C4 Project
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
	// LastCycle is the last announced cycle in this session.
	LastCycle director.Cycle `json:"last_cycle,omitempty"`

	// Session contains statistics for this session.
	Session MachineSpan `json:"session,omitempty"`
	// Total contains statistics across all sessions.
	Total MachineSpan `json:"total,omitempty"`
}

// ResetForSession removes from this statset any statistics that no longer apply across session boundaries.
func (m *Machine) ResetForSession() {
	m.LastCycle = director.Cycle{}
	m.Session.Reset()
}

// AddCycle adds the information from cycle message c to this machine statset.
func (m *Machine) AddCycle(c director.CycleMessage) {
	if c.Kind == director.CycleStart {
		m.LastCycle = c.Cycle
	}
	m.Session.AddCycle(c)
	m.Total.AddCycle(c)
}

// AddAnalysis adds the information from analysis a to this machine statset.
func (m *Machine) AddAnalysis(a analysis.Analysis) {
	m.Session.AddAnalysis(a)
	m.Total.AddAnalysis(a)
}

// MachineSpan contains the timespan-specific part of Machine.
type MachineSpan struct {
	// FinishedCycles counts the number of cycles that finished.
	FinishedCycles uint64 `json:"finished_cycles"`
	// ErroredCycles counts the number of cycles that resulted in an error.
	ErroredCycles uint64 `json:"errored_cycles"`
	// Mutation contains totals for mutation testing since this span started.
	Mutation mutation.Statset `json:"mutation,omitempty"`

	// SessionStatusTotals contains status totals since this span started.
	// It may be empty if this machine has not yet been active this span.
	StatusTotals map[status.Status]uint64 `json:"status_totals,omitempty"`
}

// Reset resets a machine span.
func (m *MachineSpan) Reset() {
	m.FinishedCycles = 0
	m.ErroredCycles = 0
	m.StatusTotals = make(map[status.Status]uint64)
	m.Mutation.Reset()
}

// AddCycle adds the information from cycle message c to this machine span.
func (m *MachineSpan) AddCycle(c director.CycleMessage) {
	switch c.Kind {
	case director.CycleFinish:
		m.FinishedCycles++
	case director.CycleError:
		m.ErroredCycles++
	}
}

// AddAnalysis adds the information from analysis a to this machine statset.
func (m *MachineSpan) AddAnalysis(a analysis.Analysis) {
	m.addStatusTotals(a)
	m.addMutation(a)
}

func (m *MachineSpan) addStatusTotals(a analysis.Analysis) {
	if m.StatusTotals == nil {
		m.StatusTotals = make(map[status.Status]uint64)
	}
	for i := status.Ok; i <= status.Last; i++ {
		count := uint64(len(a.ByStatus[i]))
		m.StatusTotals[i] += count
	}
}

func (m *MachineSpan) addMutation(a analysis.Analysis) {
	if len(a.Mutation) == 0 {
		return
	}
	m.Mutation.AddAnalysis(a.Mutation)
}
