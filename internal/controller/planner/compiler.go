// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/helper/stringhelp"

	"github.com/MattWindsor91/act-tester/internal/model/compiler/optlevel"

	"github.com/MattWindsor91/act-tester/internal/model/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// CompilerObserver observes the actions of a CompilerPlanner.
type CompilerObserver interface {
	// OnCompilerPlanStart observes that the compiler planner is beginning to configure ncompilers compilers.
	OnCompilerPlanStart(ncompilers int)
	// OnCompilerPlan observes that the corpus has added the compiler c to the plan.
	OnCompilerPlan(c compiler.Named)
	// OnCompilerPlanFinish observes that the compiler planner has finished adding compilers.
	OnCompilerPlanFinish()
}

// OnCompilerPlanStart sends an OnCompilerPlanStart to every observer in obs.
func OnCompilerPlanStart(ncompilers int, obs ...CompilerObserver) {
	for _, o := range obs {
		o.OnCompilerPlanStart(ncompilers)
	}
}

// OnCompilerPlan sends an OnCompilerPlanStart to every observer in obs.
func OnCompilerPlan(c compiler.Named, obs ...CompilerObserver) {
	for _, o := range obs {
		o.OnCompilerPlan(c)
	}
}

// OnCompilerPlanFinish sends an OnCompilerPlanStart to every observer in obs.
func OnCompilerPlanFinish(obs ...CompilerObserver) {
	for _, o := range obs {
		o.OnCompilerPlanFinish()
	}
}

// CompilerLister is the interface of things that can query compiler information.
type CompilerLister interface {
	// ListCompilers asks the compiler inspector to list all available compilers on machine ID mid.
	ListCompilers(ctx context.Context, mid id.ID) (map[string]compiler.Config, error)
}

// CompilerPlanner contains the state necessary to make up the compiler part of a test plan.
type CompilerPlanner struct {
	// Lister lists the available compilers.
	Lister CompilerLister
	// Inspector resolves configuration pertaining to a particular compiler.
	Inspector compiler.Inspector
	// Observers contains observers for the CompilerPlanner.
	Observers []CompilerObserver
	// MachineID is the identifier of the machine for which we are making a plan.
	MachineID id.ID
	// Rng is the random number generator to use in configuration randomisation.
	Rng *rand.Rand
}

func (p *Planner) planCompilers(ctx context.Context) error {
	c := CompilerPlanner{
		Lister:    p.conf.Source.CLister,
		Inspector: p.conf.Source.CInspector,
		Observers: p.conf.Observers.Compiler,
		MachineID: p.plan.Machine.ID,
		Rng:       p.rng,
	}
	var err error
	p.plan.Compilers, err = c.Plan(ctx)
	return err
}

// Plan constructs the compiler set for a plan.
func (c *CompilerPlanner) Plan(ctx context.Context) (map[string]compiler.Compiler, error) {
	cfgs, err := c.Lister.ListCompilers(ctx, c.MachineID)
	if err != nil {
		return nil, fmt.Errorf("listing compilers: %w", err)
	}

	nenabled := resolveDisabled(cfgs)
	OnCompilerPlanStart(nenabled, c.Observers...)

	cmps := make(map[string]compiler.Compiler, len(cfgs))
	for n, cfg := range cfgs {
		nc, err := c.maybePlanCompiler(cmps, n, cfg)
		if err != nil {
			return nil, err
		}
		if nc != nil {
			OnCompilerPlan(*nc, c.Observers...)
		}
	}

	OnCompilerPlanFinish(c.Observers...)

	return cmps, nil
}

func resolveDisabled(cfgs map[string]compiler.Config) (nenabled int) {
	// TODO(@MattWindsor91): automatic disabling
	for _, cfg := range cfgs {
		if !cfg.Disabled {
			nenabled++
		}
	}
	return nenabled
}

func (c *CompilerPlanner) maybePlanCompiler(into map[string]compiler.Compiler, n string, cfg compiler.Config) (*compiler.Named, error) {
	if cfg.Disabled {
		return nil, nil
	}

	nid, err := id.TryFromString(n)
	if err != nil {
		return nil, fmt.Errorf("%s not a valid ID: %w", n, err)
	}
	if into[n], err = c.planCompiler(cfg); err != nil {
		return nil, fmt.Errorf("planning compiler %s/%s: %w", c.MachineID, n, err)
	}

	return into[n].AddName(nid), nil
}

func (c *CompilerPlanner) planCompiler(cfg compiler.Config) (compiler.Compiler, error) {
	opt, err := c.planCompilerOpt(cfg)
	if err != nil {
		return compiler.Compiler{}, nil
	}
	mopt, err := c.planCompilerMOpt(cfg)
	comp := compiler.Compiler{
		SelectedOpt:  opt,
		SelectedMOpt: mopt,
		Config:       cfg,
	}
	return comp, err
}

func (c *CompilerPlanner) planCompilerOpt(cfg compiler.Config) (*optlevel.Named, error) {
	opts, err := compiler.SelectLevels(c.Inspector, &cfg)
	if err != nil {
		return nil, err
	}
	names, err := stringhelp.MapKeys(opts)
	if err != nil {
		return nil, err
	}
	return c.chooseOpt(opts, names), err
}

func (c *CompilerPlanner) planCompilerMOpt(cfg compiler.Config) (string, error) {
	mopts, err := compiler.SelectMOpts(c.Inspector, &cfg)
	if err != nil {
		return "", err
	}
	return c.chooseMOpt(mopts), err
}

func (c *CompilerPlanner) chooseOpt(opts map[string]optlevel.Level, names []string) *optlevel.Named {
	// Don't bother trying to select an optimisation if there aren't any
	if len(opts) == 0 {
		return nil
	}

	// The idea here is that we're giving 'don't choose an optimisation' - index -1 - an equal chance.
	i := c.Rng.Intn(len(opts)+1) - 1
	if i < 0 {
		return nil
	}

	name := names[i]
	return &optlevel.Named{Name: name, Level: opts[name]}

}

func (c *CompilerPlanner) chooseMOpt(opts stringhelp.Set) string {
	// Don't bother trying to select an mopt if there aren't any
	// TODO(@MattWindsor91): should this be an error?
	nopts := len(opts)
	if nopts == 0 {
		return ""
	}
	// 'don't choose an mopt' - the empty string, may or may not be a valid choice, so we don't factor it in here.
	optsl := opts.Slice()
	i := c.Rng.Intn(nopts)
	return optsl[i]
}
