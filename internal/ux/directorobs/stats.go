// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package directorobs

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/MattWindsor91/c4t/internal/director/pathset"
	"github.com/MattWindsor91/c4t/internal/quantity"

	"github.com/MattWindsor91/c4t/internal/machine"

	"github.com/MattWindsor91/c4t/internal/helper/iohelp"
	"github.com/MattWindsor91/c4t/internal/plan/analysis"
	"github.com/MattWindsor91/c4t/internal/subject/status"

	"github.com/MattWindsor91/c4t/internal/helper/errhelp"

	"github.com/MattWindsor91/c4t/internal/director"
	"github.com/MattWindsor91/c4t/internal/model/service/compiler"
	"github.com/MattWindsor91/c4t/internal/stage/analyser/saver"
)

// Statset aggregates statistics taken from cycle analyses.
type Statset struct {
	// TODO(@MattWindsor91): use this to drive the dashboard, eg by generalising the persister and giving it hints
	// as to which machine has changed

	// StartTime is the time at which this statistics set was first created.
	StartTime time.Time `json:"start_time"`

	// EventCount is the number of events that have been applied to this Statset.
	// This metric isn't particularly exciting from a user perspective, but is used to prevent spurious disk flushes.
	EventCount uint64 `json:"event_count"`

	// SessionStartTime is the time at which this session started (ie, the statset was last reloaded from disk).
	SessionStartTime time.Time `json:"session_start_time"`

	// Machines is a map from machine IDs to statistics about those machines.
	Machines map[string]MachineStatset `json:"machines,omitempty"`
}

// OnCycle incorporates cycle information from c into the statistics set.
func (s *Statset) OnCycle(c director.CycleMessage) {
	s.liftCycle(c.Cycle, func(m *MachineStatset) {
		if c.Kind == director.CycleStart {
			m.TotalCycles++
			m.LastCycle = c.Cycle
			s.EventCount++
		}
	})
}

// OnAnalysis incorporates cycle analysis from a into the statistics set.
func (s *Statset) OnCycleAnalysis(a director.CycleAnalysis) {
	s.liftCycle(a.Cycle, func(m *MachineStatset) {
		m.AddAnalysis(a.Analysis)
		s.EventCount++
	})
}

// OnCycleCompiler does nothing, for now.
func (s *Statset) OnCycleCompiler(director.Cycle, compiler.Message) {
}

// OnCycleSave does nothing, for now.
func (s *Statset) OnCycleSave(director.Cycle, saver.ArchiveMessage) {
}

// OnMachines does nothing, for now.
func (s *Statset) OnMachines(message machine.Message) {
}

// OnPrepare does nothing, for now.
func (s *Statset) OnPrepare(quantity.RootSet, pathset.Pathset) {
}

// liftCycle lifts an operation on a cycle c, handling making sure the machine exists, translating the ID, etc.
func (s *Statset) liftCycle(c director.Cycle, f func(m *MachineStatset)) {
	if s.Machines == nil {
		s.Machines = make(map[string]MachineStatset)
	}
	mid := c.MachineID.String()
	mach := s.Machines[mid]
	f(&mach)
	s.Machines[mid] = mach
}

// ResetForSession resets statistics in s that are session-specific.
func (s *Statset) ResetForSession() {
	s.SessionStartTime = time.Now()
	for k, mach := range s.Machines {
		mach.ResetForSession()
		s.Machines[k] = mach
	}
}

// Init initialises statistics in s that should be set when creating a stats file.
func (s *Statset) Init() {
	s.StartTime = time.Now()
}

// MachineStatset is a statistics set for a specific machine.
type MachineStatset struct {
	// TotalCycles is the number of the cycles the machine has accrued.
	TotalCycles uint64 `json:"total_cycles"`
	// LastCycle is the last announced cycle in this session.
	LastCycle director.Cycle `json:"last_cycle,omitempty"`

	// StatusTotals contains status totals since this statset started.
	StatusTotals map[status.Status]uint64 `json:"status_totals"`
	// SessionStatusTotals contains status totals since this session started.
	SessionStatusTotals map[status.Status]uint64 `json:"session_status_totals"`
}

// ResetForSession removes from this statset any statistics that no longer apply across session boundaries.
func (m *MachineStatset) ResetForSession() {
	m.LastCycle = director.Cycle{}
	m.SessionStatusTotals = make(map[status.Status]uint64)
}

// AddAnalysis adds the information from analysis a to this machine statset.
func (m *MachineStatset) AddAnalysis(a analysis.Analysis) {
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

// StatPersister is a forward handler that maintains and persists a statistics set on disk.
type StatPersister struct {
	// set is the statistics set being persisted.
	set Statset
	// f is the target file (we need a file to be able to truncate properly).
	f *os.File
	// enc is the encoder writing to f.
	enc *json.Encoder
	// err holds any error that has occurred while persisting statistics.
	err error
	// lastCount is the value of EventCount when set was last committed to disk.
	lastCount uint64
}

// NewStatPersister creates a StatPersister that reads and writes statistics from f.
// If f is non-empty, it immediately tries to read any existing stats dump, and fails if this doesn't work.
// The StatPersister takes ownership of f; close f with Close.
func NewStatPersister(f *os.File) (*StatPersister, error) {
	sp := StatPersister{f: f, enc: json.NewEncoder(f)}
	if err := sp.tryReadStats(); err != nil {
		return nil, fmt.Errorf("while reloading stats: %w", err)
	}
	sp.set.ResetForSession()
	return &sp, nil
}

// Close closes this StatPersister, returning any errors arising from either the stats persisting or file close.
func (s *StatPersister) Close() error {
	perr := s.err
	cerr := s.f.Close()
	return errhelp.FirstError(perr, cerr)
}

// OnMachine feeds the information from m into the stats set.
func (s *StatPersister) OnMachines(m machine.Message) {
	s.set.OnMachines(m)
	s.flush()
}

// OnMachine feeds the information from m into the stats set.
func (s *StatPersister) OnPrepare(qs quantity.RootSet, ps pathset.Pathset) {
	s.set.OnPrepare(qs, ps)
	s.flush()
}

// OnCycle feeds the information from c into the stats set.
func (s *StatPersister) OnCycle(c director.CycleMessage) {
	s.set.OnCycle(c)
	s.flush()
}

// OnCycleAnalysis feeds the information from a into the stats set.
func (s *StatPersister) OnCycleAnalysis(a director.CycleAnalysis) {
	s.set.OnCycleAnalysis(a)
	s.flush()
}

// OnCycleCompiler feeds the information from c and m into the stats set.
func (s *StatPersister) OnCycleCompiler(c director.Cycle, m compiler.Message) {
	s.set.OnCycleCompiler(c, m)
	s.flush()
}

// OnCycleSave feeds the information from c and m into the stats set.
func (s *StatPersister) OnCycleSave(c director.Cycle, m saver.ArchiveMessage) {
	s.set.OnCycleSave(c, m)
	s.flush()
}

func (s *StatPersister) tryReadStats() error {
	if empty, err := iohelp.IsFileEmpty(s.f); err != nil {
		return fmt.Errorf("while checking file for existing stats: %w", err)
	} else if empty {
		s.set.StartTime = time.Now()
		return nil
	}
	return s.set.Load(s.f)
}

// Load loads stats from r into this Statset.
func (s *Statset) Load(r io.Reader) error {
	return json.NewDecoder(r).Decode(s)
}

func (s *StatPersister) flush() {
	if s.err != nil {
		return
	}
	// Possible ABA problem, but should be ok as long as at most one event happens before each flush.
	if s.set.EventCount == s.lastCount {
		return
	}
	s.lastCount = s.set.EventCount
	if _, s.err = s.f.Seek(0, io.SeekStart); s.err != nil {
		return
	}
	if s.err = s.f.Truncate(0); s.err != nil {
		return
	}
	s.err = s.enc.Encode(s.set)
}

// OpenStatFile opens a file in the appropriate mode for using it as a statistics target.
func OpenStatFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
}
