// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package quantity

import (
	"log"
)

// RootSet is the top-level set of tunable quantities for a director.
type RootSet struct {
	// GlobalTimeout is the top-level timeout for the director.
	GlobalTimeout Timeout `toml:"global_timeout,omitzero"`

	// Plan is the quantity set for the planner stage.
	Plan PlanSet `toml:"plan,omitempty"`

	// This part of the quantity set is effectively a default for all machines that don't have overrides.
	MachineSet
}

// Log logs q to l.
func (q *RootSet) Log(l *log.Logger) {
	l.Println("[Plan]")
	q.Plan.Log(l)
	q.MachineSet.Log(l)
}

// Override substitutes any quantities in new that are non-zero for those in this set.
func (q *RootSet) Override(new RootSet) {
	q.GlobalTimeout.Override(new.GlobalTimeout)
	q.Plan.Override(new.Plan)
	q.MachineSet.Override(new.MachineSet)
}
