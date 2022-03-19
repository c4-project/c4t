// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"fmt"

	"github.com/c4-project/c4t/internal/machine"

	"github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/c4-project/c4t/internal/id"
)

// CompilerLister is the interface of things that can query compiler information for a particular machine.
type CompilerLister interface {
	// Compilers asks the compiler inspector to list all available compilers.
	Compilers() (map[id.ID]compiler.Compiler, error)
}

//go:generate mockery --name=CompilerLister

// CompilerPlanner contains the state necessary to make up the compiler part of a test plan.
type CompilerPlanner struct {
	// Lister lists the available compilers.
	Lister CompilerLister
	// Filter is the filtering glob to use on compiler names.
	Filter id.ID
	// Observers contains observers for the CompilerPlanner.
	Observers []compiler.Observer
}

func (p *Planner) planCompilers(m machine.Config) (compiler.InstanceMap, error) {
	c := CompilerPlanner{
		Filter:    id.FromString(p.filter),
		Observers: lowerToCompiler(p.observers),
		Lister:    &m,
	}
	return c.Plan()
}

// Plan constructs the compiler set for a plan.
func (c *CompilerPlanner) Plan() (compiler.InstanceMap, error) {
	cfgs, err := c.Lister.Compilers()
	if err != nil {
		return nil, fmt.Errorf("listing compilers: %w", err)
	}

	if cfgs, err = c.filterCompilers(cfgs); err != nil {
		return nil, fmt.Errorf("filtering compilers: %w", err)
	}

	nenabled := resolveDisabled(cfgs)
	compiler.OnCompilerConfigStart(nenabled, c.Observers...)

	cmps := make(compiler.InstanceMap, len(cfgs))
	i := 0
	for n, cfg := range cfgs {
		nc := c.maybePlanCompiler(cmps, n, cfg)
		if nc != nil {
			compiler.OnCompilerConfigStep(i, *nc, c.Observers...)
		}
		i++
	}

	compiler.OnCompilerConfigEnd(c.Observers...)

	return cmps, nil
}

func (c *CompilerPlanner) filterCompilers(in map[id.ID]compiler.Compiler) (map[id.ID]compiler.Compiler, error) {
	if c.Filter.IsEmpty() {
		return in, nil
	}
	out, err := id.MapGlob(in, c.Filter)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func resolveDisabled(cfgs map[id.ID]compiler.Compiler) (nenabled int) {
	// TODO(@MattWindsor91): automatic disabling
	for _, cfg := range cfgs {
		if !cfg.Disabled {
			nenabled++
		}
	}
	return nenabled
}

func (c *CompilerPlanner) maybePlanCompiler(into compiler.InstanceMap, nid id.ID, cfg compiler.Compiler) *compiler.Named {
	if cfg.Disabled {
		return nil
	}
	// Everything that used to be here is now in the perturber.
	into[nid] = compiler.Instance{Compiler: cfg}
	return into[nid].AddName(nid)
}
