// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package perturber

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/plan"

	"github.com/MattWindsor91/act-tester/internal/helper/stringhelp"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler/optlevel"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// CompilerLister is the interface of things that can query compiler information.
type CompilerLister interface {
	// ListCompilers asks the compiler inspector to list all available compilers on machine ID mid.
	ListCompilers(ctx context.Context, mid id.ID) (map[string]compiler.Compiler, error)
}

// CompilerPerturber contains the state necessary to make up the compiler part of a test plan.
type CompilerPerturber struct {
	// Inspector resolves configuration pertaining to a particular compiler.
	Inspector compiler.Inspector
	// Observers contains observers for the CompilerPlanner.
	Observers []compiler.Observer
	// Rng is the random number generator to use in configuration randomisation.
	Rng *rand.Rand
}

func (p *Perturber) perturbCompilers(rng *rand.Rand, pn *plan.Plan) error {
	c := CompilerPerturber{
		Inspector: p.ci,
		Observers: lowerToCompiler(p.observers),
		Rng:       rng,
	}
	var err error
	pn.Compilers, err = c.Perturb(pn.Compilers)
	return err
}

func lowerToCompiler(obs []Observer) []compiler.Observer {
	cobs := make([]compiler.Observer, len(obs))
	for i, o := range obs {
		cobs[i] = o
	}
	return cobs
}

// Perturb perturbs the compiler set for a plan.
func (c *CompilerPerturber) Perturb(cfgs map[string]compiler.Configuration) (map[string]compiler.Configuration, error) {
	compiler.OnCompilerConfigStart(len(cfgs), c.Observers...)

	ncfgs := make(map[string]compiler.Configuration, len(cfgs))
	i := 0
	for n, cfg := range cfgs {
		nc, err := c.perturbCompiler(n, cfg.Compiler)
		if err != nil {
			return nil, err
		}
		ncfgs[n] = nc.Configuration
		compiler.OnCompilerConfigStep(i, *nc, c.Observers...)
		i++
	}

	compiler.OnCompilerConfigEnd(c.Observers...)

	return ncfgs, nil
}

func (c *CompilerPerturber) perturbCompiler(name string, cmp compiler.Compiler) (*compiler.Named, error) {
	nid, err := id.TryFromString(name)
	if err != nil {
		return nil, fmt.Errorf("%s not a valid ID: %w", name, err)
	}

	opt, err := c.perturbCompilerOpt(cmp)
	if err != nil {
		return nil, err
	}
	mopt, err := c.perturbCompilerMOpt(cmp)
	comp := compiler.Configuration{
		SelectedOpt:  opt,
		SelectedMOpt: mopt,
		Compiler:     cmp,
	}

	return comp.AddName(nid), err
}

func (c *CompilerPerturber) perturbCompilerOpt(cfg compiler.Compiler) (*optlevel.Named, error) {
	opts, err := compiler.SelectLevels(c.Inspector, &cfg)
	if err != nil {
		return nil, err
	}
	names, err := stringhelp.MapKeys(opts)
	if err != nil {
		return nil, err
	}
	return c.chooseOpt(opts, names), nil
}

func (c *CompilerPerturber) perturbCompilerMOpt(cfg compiler.Compiler) (string, error) {
	mopts, err := compiler.SelectMOpts(c.Inspector, &cfg)
	if err != nil {
		return "", err
	}
	return c.chooseMOpt(mopts), nil
}

func (c *CompilerPerturber) chooseOpt(opts map[string]optlevel.Level, names []string) *optlevel.Named {
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

func (c *CompilerPerturber) chooseMOpt(opts stringhelp.Set) string {
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
