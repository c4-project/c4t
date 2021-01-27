// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stat

import (
	"encoding/json"
	"io"
	"time"

	"github.com/c4-project/c4t/internal/copier"
	"github.com/c4-project/c4t/internal/director"
	"github.com/c4-project/c4t/internal/machine"
	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/stage/analyser/saver"
	"github.com/c4-project/c4t/internal/subject/corpus/builder"
)

// Set aggregates statistics taken from cycle analyses.
type Set struct {
	// TODO(@MattWindsor91): use this to drive the dashboard, eg by generalising the persister and giving it hints
	// as to which machine has changed

	// StartTime is the time at which this statistics set was first created.
	StartTime time.Time `json:"start_time,omitempty"`

	// EventCount is the number of events that have been applied to this Set.
	// This metric isn't particularly exciting from a user perspective, but is used to prevent spurious disk flushes.
	EventCount uint64 `json:"event_count"`

	// SessionStartTime is the time at which this session started (ie, the statset was last reloaded from disk).
	SessionStartTime time.Time `json:"session_start_time,omitempty"`

	// Machines is a map from machine IDs to statistics about those machines.
	Machines map[string]Machine `json:"machines,omitempty"`
}

// OnCycle incorporates cycle information from c into the statistics set.
func (s *Set) OnCycle(c director.CycleMessage) {
	s.liftCycle(c.Cycle, func(m *Machine) {
		m.AddCycle(c)
	})
	s.EventCount++
}

// OnAnalysis incorporates cycle analysis from a into the statistics set.
func (s *Set) OnCycleAnalysis(a director.CycleAnalysis) {
	s.liftCycle(a.Cycle, func(m *Machine) {
		m.AddAnalysis(a.Analysis)
	})
	s.EventCount++
}

// OnCycleBuild does nothing, for now.
func (s *Set) OnCycleBuild(director.Cycle, builder.Message) {}

// OnCycleCompiler does nothing, for now.
func (s *Set) OnCycleCompiler(director.Cycle, compiler.Message) {}

// OnCycleCopy does nothing, for now.
func (s *Set) OnCycleCopy(director.Cycle, copier.Message) {}

// OnCycleInstance does nothing, for now.
func (s *Set) OnCycleInstance(director.Cycle, director.InstanceMessage) {}

// OnCycleSave does nothing, for now.
func (s *Set) OnCycleSave(director.Cycle, saver.ArchiveMessage) {}

// OnMachines does nothing, for now.
func (s *Set) OnMachines(machine.Message) {}

// OnPrepare does nothing, for now.
func (s *Set) OnPrepare(director.PrepareMessage) {}

// liftCycle lifts an operation on a cycle c, handling making sure the machine exists, translating the ID, etc.
func (s *Set) liftCycle(c director.Cycle, f func(m *Machine)) {
	if s.Machines == nil {
		s.Machines = make(map[string]Machine)
	}
	mid := c.MachineID.String()
	mach := s.Machines[mid]
	f(&mach)
	s.Machines[mid] = mach
}

// ResetForSession resets statistics in s that are session-specific.
func (s *Set) ResetForSession() {
	s.SessionStartTime = time.Now()
	for k, mach := range s.Machines {
		mach.ResetForSession()
		s.Machines[k] = mach
	}
}

// Init initialises statistics in s that should be set when creating a stats file.
func (s *Set) Init() {
	s.StartTime = time.Now()
}

// Load loads stats from r into this Set.
func (s *Set) Load(r io.Reader) error {
	return json.NewDecoder(r).Decode(s)
}
