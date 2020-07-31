// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"context"
	"errors"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/stage/perturber"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser"

	"github.com/MattWindsor91/act-tester/internal/stage/invoker"

	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/stage/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/stage/lifter"
	"github.com/MattWindsor91/act-tester/internal/stage/planner"
)

// StageConfig groups together the stage configuration for of a director instance.
type StageConfig struct {
	// Plan contains configuration for the instance's plan stage.
	// TODO(@MattWindsor91): remove this
	Plan *planner.Planner
	// Perturb contains configuration for the instance's perturb stage.
	Perturb *perturber.Perturber
	// Fuzz contains configuration for the instance's fuzz stage.
	Fuzz *fuzzer.Fuzzer
	// Lift contains configuration for the instance's lift stage.
	Lift *lifter.Lifter
	// Invoke contains configuration for the instance's invoke stage.
	Invoke *invoker.Invoker
	// Analyser contains configuration for the instance's analyser stage.
	Analyser *analyser.Analyser
}

var ErrStageConfigMissing = errors.New("stage config missing")

// Check makes sure the StageConfig has all configuration elements present.
func (c *StageConfig) Check() error {
	if c.Plan == nil {
		return fmt.Errorf("%w: %s", ErrStageConfigMissing, stage.Plan)
	}
	if c.Perturb == nil {
		return fmt.Errorf("%w: %s", ErrStageConfigMissing, stage.Perturb)
	}
	if c.Fuzz == nil {
		return fmt.Errorf("%w: %s", ErrStageConfigMissing, stage.Fuzz)
	}
	if c.Lift == nil {
		return fmt.Errorf("%w: %s", ErrStageConfigMissing, stage.Lift)
	}
	if c.Invoke == nil {
		return fmt.Errorf("%w: %s", ErrStageConfigMissing, stage.Invoke)
	}
	if c.Analyser == nil {
		return fmt.Errorf("%w: %s", ErrStageConfigMissing, stage.Analyse)
	}
	return nil
}

// stageRunner is the type of encapsulated director stages.
type stageRunner struct {
	// Stage is the ID of the stage, which appears in logging and errors.
	Stage stage.Stage
	// Run is the function to use to run the stage.
	Run func(*StageConfig, context.Context, *plan.Plan) (*plan.Plan, error)
}

// Stages is the list of director stages.
var Stages = []stageRunner{
	{
		Stage: stage.Plan,
		Run: func(c *StageConfig, ctx context.Context, _ *plan.Plan) (*plan.Plan, error) {
			return c.Plan.Plan(ctx)
		},
	},
	{
		Stage: stage.Perturb,
		Run: func(c *StageConfig, ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
			return c.Perturb.Run(ctx, p)
		},
	},
	{
		Stage: stage.Fuzz,
		Run: func(c *StageConfig, ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
			return c.Fuzz.Run(ctx, p)
		},
	},
	{
		Stage: stage.Lift,
		Run: func(c *StageConfig, ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
			return c.Lift.Run(ctx, p)
		},
	},
	{
		Stage: stage.Invoke,
		Run: func(c *StageConfig, ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
			return c.Invoke.Run(ctx, p)
		},
	},
	{
		Stage: stage.Analyse,
		Run: func(c *StageConfig, ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
			return c.Analyser.Run(ctx, p)
		},
	},
}
