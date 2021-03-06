// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"fmt"

	"github.com/c4-project/c4t/internal/director/pathset"
	"github.com/c4-project/c4t/internal/helper/iohelp"
	"github.com/c4-project/c4t/internal/id"
	"github.com/c4-project/c4t/internal/machine"
	"github.com/c4-project/c4t/internal/plan"
	"github.com/c4-project/c4t/internal/quantity"
)

// TODO(@MattWindsor91): make a proper distinction between instances and machines.

// Machine contains the state for a particular machine attached to an instance.
type Machine struct {
	// ID is the ID for this machine.
	ID id.ID

	// InitialPlan is the plan that is perturbed to form the plan for each test cycle.
	InitialPlan plan.Plan

	// Pathset contains the pathset for this instance.
	Pathset *pathset.Instance

	// Quantities contains the quantity set for this machine.
	Quantities quantity.MachineSet

	// Config contains the machine config for this machine.
	Config machine.Config

	// cycle is the number of the current cycle for the machine.
	// This is held separately from the instance as an instance may (eventually) run cycles for multiple machines.
	cycle uint64

	// stages is the set of stages to run for this machine.
	stages []plan.Runner
}

func (m *Machine) check() error {
	if m.Pathset == nil {
		return fmt.Errorf("%w: paths for machine %s", iohelp.ErrPathsetNil, m.ID.String())
	}
	return nil
}
