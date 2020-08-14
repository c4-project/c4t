// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package quantity

import (
	"log"
)

// MachineSet contains overridable quantities for each stage operating on a particular machine.
// Often, but not always, these quantities will be shared between machines.
type MachineSet struct {
	// Fuzz is the quantity set for the fuzz stage.
	Fuzz FuzzSet `toml:"fuzz,omitzero"`
	// Mach is the quantity set for the machine-local stage, as well as any machine-local stages run remotely.
	Mach MachNodeSet `toml:"mach,omitzero"`
	// Perturb is the quantity set for the planner stage.
	Perturb PerturbSet `toml:"perturb,omitzero"`
}

// Log logs q to l.
func (q *MachineSet) Log(l *log.Logger) {
	l.Println("[Perturb]")
	q.Perturb.Log(l)
	l.Println("[Fuzz]")
	q.Fuzz.Log(l)
	l.Println("[Mach]")
	q.Mach.Log(l)
}

// Override substitutes any quantities in new that are non-zero for those in this set.
func (q *MachineSet) Override(new MachineSet) {
	q.Perturb.Override(new.Perturb)
	q.Fuzz.Override(new.Fuzz)
	q.Mach.Override(new.Mach)
}
