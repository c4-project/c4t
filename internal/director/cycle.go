// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"context"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/run"
	"github.com/MattWindsor91/act-tester/internal/plan"
)

// cycle contains state for one cycle of one instance of the director.
type cycle struct {
	// header contains information about which cycle this is.
	header run.Run

	// p points to the plan being built in this cycle.
	p *plan.Plan

	// sc points to the instance's stage config.
	sc *StageConfig
}

func (c *cycle) run(ctx context.Context) error {
	for _, s := range Stages {
		if err := c.runStage(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func (c *cycle) runStage(ctx context.Context, s stageRunner) error {
	var err error
	if c.p, err = s.Run(c.sc, ctx, c.p); err != nil {
		return fmt.Errorf("in %s stage: %w", s.Stage, err)
	}
	return nil
}
