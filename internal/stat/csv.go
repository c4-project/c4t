// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stat

import (
	"encoding/csv"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/subject/status"
)

// DumpMutationCSVHeader dumps into w a CSV header for mutation analysis.
func (s *Set) DumpMutationCSVHeader(w *csv.Writer) error {
	hdr := []string{
		"Machine",
		"Index",
		"Name",
		"Selections",
		"Hits",
		"Kills",
	}
	for i := status.Ok; i <= status.Last; i++ {
		hdr = append(hdr, i.String())
	}
	return w.Write(hdr)
}

// DumpMutationCSV dumps into w a CSV representation of the mutation statistics in this set.
// Each machine record has its lines prefixed by its machine ID, is flushed separately, and appears in ID order.
// If total is true, the multi-session totals will be dumped; otherwise, this session's totals will be dumped.
func (s *Set) DumpMutationCSV(w *csv.Writer, total bool) error {
	for _, mid := range id.MapKeys(s.Machines) {
		sm := s.Machines[mid]
		if err := sm.DumpMutationCSV(w, mid, total); err != nil {
			return err
		}
	}
	return nil
}

// DumpMutationCSV dumps into w a CSV representation of the mutation statistics in this machine.
// Each line in the record has mid as a prefix.
// If total is true, the multi-session totals will be dumped; otherwise, this session's totals will be dumped.
// The writer is flushed at the end of this dump.
func (m *Machine) DumpMutationCSV(w *csv.Writer, machine id.ID, total bool) error {
	if total {
		return m.Total.DumpMutationCSV(w, machine)
	}
	return m.Session.DumpMutationCSV(w, machine)
}

// DumpMutationCSV dumps into w a CSV representation of the mutation statistics in this machine span.
// Each line in the record has mid as a prefix.
// The writer is flushed at the end of this dump.
func (m *MachineSpan) DumpMutationCSV(w *csv.Writer, machine id.ID) error {
	return m.Mutation.DumpCSV(w, machine)
}
