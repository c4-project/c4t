// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package quantity

import (
	"log"
)

// RootSet is the top-level set of tunable quantities for a director.
type RootSet struct {
	// Plan is the quantity set for the planner stage.
	Plan PlanSet `toml:"plan,omitzero"`

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
	q.Plan.Override(new.Plan)
	q.MachineSet.Override(new.MachineSet)
}
