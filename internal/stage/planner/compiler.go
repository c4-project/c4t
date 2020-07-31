// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"context"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/plan"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// CompilerLister is the interface of things that can query compiler information.
type CompilerLister interface {
	// ListCompilers asks the compiler inspector to list all available compilers on machine ID mid.
	ListCompilers(ctx context.Context, mid id.ID) (map[string]compiler.Compiler, error)
}

// CompilerPlanner contains the state necessary to make up the compiler part of a test plan.
type CompilerPlanner struct {
	// Lister lists the available compilers.
	Lister CompilerLister
	// Filter is the filtering glob to use on compiler names.
	Filter id.ID
	// Observers contains observers for the CompilerPlanner.
	Observers []compiler.Observer
	// MachineID is the identifier of the machine for which we are making a plan.
	MachineID id.ID
}

func (p *Planner) planCompilers(ctx context.Context, pn *plan.Plan) error {
	c := CompilerPlanner{
		Filter:    id.FromString(p.filter),
		Lister:    p.source.CLister,
		Observers: lowerToCompiler(p.observers),
		MachineID: pn.Machine.ID,
	}
	var err error
	pn.Compilers, err = c.Plan(ctx)
	return err
}

// Plan constructs the compiler set for a plan.
func (c *CompilerPlanner) Plan(ctx context.Context) (map[string]compiler.Configuration, error) {
	cfgs, err := c.Lister.ListCompilers(ctx, c.MachineID)
	if err != nil {
		return nil, fmt.Errorf("listing compilers: %w", err)
	}

	if cfgs, err = c.filterCompilers(cfgs); err != nil {
		return nil, fmt.Errorf("filtering compilers: %w", err)
	}

	nenabled := resolveDisabled(cfgs)
	compiler.OnCompilerConfigStart(nenabled, c.Observers...)

	cmps := make(map[string]compiler.Configuration, len(cfgs))
	i := 0
	for n, cfg := range cfgs {
		nc, err := c.maybePlanCompiler(cmps, n, cfg)
		if err != nil {
			return nil, err
		}
		if nc != nil {
			compiler.OnCompilerConfigStep(i, *nc, c.Observers...)
		}
		i++
	}

	compiler.OnCompilerConfigEnd(c.Observers...)

	return cmps, nil
}

func (c *CompilerPlanner) filterCompilers(in map[string]compiler.Compiler) (map[string]compiler.Compiler, error) {
	if c.Filter.IsEmpty() {
		return in, nil
	}
	out, err := id.MapGlob(in, c.Filter)
	if err != nil {
		return nil, err
	}
	return out.(map[string]compiler.Compiler), nil
}

func resolveDisabled(cfgs map[string]compiler.Compiler) (nenabled int) {
	// TODO(@MattWindsor91): automatic disabling
	for _, cfg := range cfgs {
		if !cfg.Disabled {
			nenabled++
		}
	}
	return nenabled
}

func (c *CompilerPlanner) maybePlanCompiler(into map[string]compiler.Configuration, n string, cfg compiler.Compiler) (*compiler.Named, error) {
	if cfg.Disabled {
		return nil, nil
	}

	nid, err := id.TryFromString(n)
	if err != nil {
		return nil, fmt.Errorf("%s not a valid ID: %w", n, err)
	}

	// Everything that used to be here is now in the perturber.
	into[n] = compiler.Configuration{Compiler: cfg}
	return into[n].AddName(nid), nil
}
