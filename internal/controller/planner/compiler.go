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
	// MachineID is the identifier of the machine for which we are making a plan.
	MachineID id.ID
	// Rng is the random number generator to use in configuration randomisation.
	Rng *rand.Rand
}

func (p *Planner) planCompilers(ctx context.Context, rng *rand.Rand) (map[string]compiler.Compiler, error) {
	c := CompilerPlanner{
		Lister:    p.Source.CLister,
		Inspector: p.Source.CInspector,
		MachineID: p.MachineID,
		Rng:       rng,
	}
	return c.Plan(ctx)
}

func (c *CompilerPlanner) Plan(ctx context.Context) (map[string]compiler.Compiler, error) {
	cfgs, err := c.Lister.ListCompilers(ctx, c.MachineID)
	if err != nil {
		return nil, fmt.Errorf("listing compilers: %w", err)
	}

	cmps := make(map[string]compiler.Compiler, len(cfgs))
	for n, cfg := range cfgs {
		var err error
		if cmps[n], err = c.planCompiler(cfg); err != nil {
			return nil, fmt.Errorf("planning compiler %s: %w", n, err)
		}
	}

	return cmps, nil
}

func (c *CompilerPlanner) planCompiler(cfg compiler.Config) (compiler.Compiler, error) {
	opt, err := c.planCompilerOpt(cfg)
	comp := compiler.Compiler{
		SelectedOpt: opt,
		Config:      cfg,
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
