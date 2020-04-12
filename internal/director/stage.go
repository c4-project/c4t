// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"context"
	"errors"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/controller/rmach"

	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/controller/lifter"
	"github.com/MattWindsor91/act-tester/internal/controller/planner"
	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// StageConfig groups together the stage configuration for of a director instance.
type StageConfig struct {
	// Plan contains configuration for the instance's plan stage.
	Plan *planner.Planner
	// Fuzz contains configuration for the instance's fuzz stage.
	Fuzz *fuzzer.Config
	// Lift contains configuration for the instance's lift stage.
	Lift *lifter.Config
	// Mach contains configuration for the instance's machine-specific stage.
	Mach *rmach.Config
	// Save contains configuration for the instance's error saving stage.
	Save *Save
}

var ErrStageConfigMissing = errors.New("stage config missing")

// Check makes sure the StageConfig has all configuration elements present.
func (c *StageConfig) Check() error {
	if c.Plan == nil {
		return fmt.Errorf("%w: %s", ErrStageConfigMissing, stagePlan)
	}
	if c.Fuzz == nil {
		return fmt.Errorf("%w: %s", ErrStageConfigMissing, stageFuzz)
	}
	if c.Lift == nil {
		return fmt.Errorf("%w: %s", ErrStageConfigMissing, stageLift)
	}
	if c.Mach == nil {
		return fmt.Errorf("%w: %s", ErrStageConfigMissing, stageMach)
	}
	if c.Save == nil {
		return fmt.Errorf("%w: %s", ErrStageConfigMissing, stageSave)
	}
	return nil
}

// Stage is the type of encapsulated director stages.
type stage struct {
	// Name is the name of the stage, which appears in logging and errors.
	Name string
	// Run is the function to use to run the stage.
	Run func(*StageConfig, context.Context, *plan.Plan) (*plan.Plan, error)
}

const (
	stagePlan = "plan"
	stageFuzz = "fuzz"
	stageLift = "lift"
	stageMach = "mach"
	stageSave = "save"
)

// Stages is the list of director stages.
var Stages = []stage{
	{
		Name: stagePlan,
		Run: func(c *StageConfig, ctx context.Context, _ *plan.Plan) (*plan.Plan, error) {
			return c.Plan.Plan(ctx)
		},
	},
	{
		Name: stageFuzz,
		Run: func(c *StageConfig, ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
			return c.Fuzz.Run(ctx, p)
		},
	},
	{
		Name: stageLift,
		Run: func(c *StageConfig, ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
			return c.Lift.Run(ctx, p)
		},
	},
	{
		Name: stageMach,
		Run: func(c *StageConfig, ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
			return c.Mach.Run(ctx, p)
		},
	},
	{
		Name: stageSave,
		Run: func(c *StageConfig, ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
			return c.Save.Run(ctx, p)
		},
	},
}
