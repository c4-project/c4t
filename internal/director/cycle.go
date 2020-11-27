// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"context"
	"fmt"
	"time"

	"github.com/MattWindsor91/c4t/internal/model/id"

	"github.com/MattWindsor91/c4t/internal/plan"
)

// cycleInstance contains state for one cycle of one instance of the director.
type cycleInstance struct {
	// cycle contains information about which cycle this is.
	cycle Cycle

	// p points to the plan being built in this cycle.
	p *plan.Plan

	// sc points to the instance's stage config.
	sc *StageConfig
}

func (c *cycleInstance) run(ctx context.Context) error {
	for _, s := range Stages {
		if err := c.runStage(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func (c *cycleInstance) runStage(ctx context.Context, s stageRunner) error {
	var err error
	if c.p, err = s.Run(c.sc, ctx, c.p); err != nil {
		return fmt.Errorf("in %s stage: %w", s.Stage, err)
	}
	return nil
}

// Cycle contains information about a particular test cycle.
type Cycle struct {
	// MachineID is the ID of the machine on which this cycle is running.
	MachineID id.ID

	// Iter is the iteration number of this cycle.
	Iter uint64

	// Start is the start time of this cycle.
	Start time.Time
}

// String returns a string containing the components of this run in a human-readable manner.
func (r Cycle) String() string {
	return fmt.Sprintf(
		"[%s #%d (%s)]",
		r.MachineID.String(),
		r.Iter,
		r.Start.Format(time.Stamp),
	)
}
