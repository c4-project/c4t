// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"context"
	"fmt"
	"time"

	"github.com/c4-project/c4t/internal/model/id"

	"github.com/c4-project/c4t/internal/plan"
)

// cycleInstance contains state for one cycle of one instance of the director.
type cycleInstance struct {
	// cycle contains information about which cycle this is.
	cycle Cycle

	// p points to the plan being built in this cycle.
	p *plan.Plan

	// stages contains the stages to run in this cycle.
	stages []plan.Runner
}

func (c *cycleInstance) run(ctx context.Context) error {
	for _, s := range c.stages {
		if err := c.runStage(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func (c *cycleInstance) runStage(ctx context.Context, s plan.Runner) error {
	var err error
	if c.p, err = c.p.RunStage(ctx, s); err != nil {
		return fmt.Errorf("in %s stage: %w", s.Stage(), err)
	}
	return nil
}

// Cycle contains information about a particular test cycle.
type Cycle struct {
	// Instance is the index of the instance running this cycle.
	// Indices currently start from 0, so the zero value is also a valid instance number; this may change.
	Instance int `json:"instance"`

	// MachineID is the ID of the machine on which this cycle is running.
	MachineID id.ID `json:"machine_id,omitempty"`

	// Iter is the iteration number of this cycle.
	// Iteration numbers currently start from 0, so the zero value is also a valid iteration number; this may change.
	Iter uint64 `json:"iter"`

	// Start is the start time of this cycle.
	Start time.Time `json:"start_time,omitempty"`
}

// String returns a string containing the components of this cycle in a human-readable manner.
func (r Cycle) String() string {
	return fmt.Sprintf(
		"[%d: %s #%d (%s)]",
		r.Instance,
		r.MachineID.String(),
		r.Iter,
		r.Start.Format(time.Stamp),
	)
}
