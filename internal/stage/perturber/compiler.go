// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package perturber

import (
	"math/rand"
	"time"

	"golang.org/x/exp/maps"

	"github.com/c4-project/c4t/internal/model/service"

	"github.com/c4-project/c4t/internal/mutation"

	"github.com/c4-project/c4t/internal/plan"

	"github.com/c4-project/c4t/internal/helper/stringhelp"

	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"

	"github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/c4-project/c4t/internal/id"
)

// compilerPerturber contains the state necessary to perturb the compiler part of a test plan.
type compilerPerturber struct {
	// inspector resolves configuration pertaining to a particular compiler.
	inspector compiler.Inspector
	// observers contains observers for the compiler perturber.
	observers []compiler.Observer
	// rng is the random number generator to use in configuration randomisation.
	rng *rand.Rand
	// useFullIDs tells the perturber whether to map compilers to full IDs.
	useFullIDs bool
	// mutant is the mutant ID to insert into compilers.
	mutant mutation.Mutant
}

func (p *Perturber) perturbCompilers(rng *rand.Rand, pn *plan.Plan) error {
	c := compilerPerturber{
		inspector:  p.ci,
		observers:  lowerToCompiler(p.observers),
		rng:        rng,
		useFullIDs: p.useFullIDs,
		mutant:     pn.Mutant(),
	}
	var err error
	pn.Compilers, err = c.Perturb(pn.Compilers)
	return err
}

// Perturb perturbs the compiler set for a plan.
func (c *compilerPerturber) Perturb(cfgs compiler.InstanceMap) (compiler.InstanceMap, error) {
	compiler.OnCompilerConfigStart(len(cfgs), c.observers...)

	ncfgs := make(compiler.InstanceMap, len(cfgs))
	i := 0
	for n, cfg := range cfgs {
		nc, err := c.perturbCompiler(n, cfg.Compiler)
		if err != nil {
			return nil, err
		}
		nid, err := c.fullCompilerName(nc)
		if err != nil {
			return nil, err
		}
		ncfgs[nid] = nc.Instance
		compiler.OnCompilerConfigStep(i, *nc, c.observers...)
		i++
	}

	compiler.OnCompilerConfigEnd(c.observers...)

	return ncfgs, nil
}

func (c *compilerPerturber) fullCompilerName(nc *compiler.Named) (id.ID, error) {
	if !c.useFullIDs {
		return nc.ID, nil
	}

	fid, err := nc.FullID()
	if err != nil {
		return id.ID{}, err
	}
	return fid, nil
}

func (c *compilerPerturber) perturbCompiler(name id.ID, cmp compiler.Compiler) (*compiler.Named, error) {
	inst, err := c.makeCompilerInstance(cmp)
	if err != nil {
		return nil, err
	}
	return inst.AddName(name), nil
}

func (c *compilerPerturber) makeCompilerInstance(cmp compiler.Compiler) (compiler.Instance, error) {
	opt, err := c.perturbCompilerOpt(cmp)
	if err != nil {
		return compiler.Instance{}, err
	}
	mopt, err := c.perturbCompilerMOpt(cmp)
	if err != nil {
		return compiler.Instance{}, err
	}
	inst := compiler.Instance{
		ConfigTime:   time.Now(),
		Mutant:       c.mutant,
		SelectedOpt:  opt,
		SelectedMOpt: mopt,
		Compiler:     cmp,
	}
	inst.Run, err = c.expandRun(inst.Run, inst.Interpolations())
	return inst, err
}

func (c *compilerPerturber) expandRun(r *service.RunInfo, interps map[string]string) (*service.RunInfo, error) {
	if r == nil {
		return r, nil
	}
	// r might point directly to the run information of the initial plan, so we need to make a copy of it.
	newr := *r
	if err := newr.Interpolate(interps); err != nil {
		return nil, err
	}
	return &newr, nil
}

func (c *compilerPerturber) perturbCompilerOpt(cfg compiler.Compiler) (*optlevel.Named, error) {
	opts, err := compiler.SelectLevels(c.inspector, &cfg)
	if err != nil {
		return nil, err
	}
	return c.chooseOpt(opts, maps.Keys(opts)), nil
}

func (c *compilerPerturber) perturbCompilerMOpt(cfg compiler.Compiler) (string, error) {
	mopts, err := compiler.SelectMOpts(c.inspector, &cfg)
	if err != nil {
		return "", err
	}
	return c.chooseMOpt(mopts), nil
}

func (c *compilerPerturber) chooseOpt(opts map[string]optlevel.Level, names []string) *optlevel.Named {
	// Don't bother trying to select an optimisation if there aren't any
	if len(opts) == 0 {
		return nil
	}

	// The idea here is that we're giving 'don't choose an optimisation' - index -1 - an equal chance.
	i := c.rng.Intn(len(opts)+1) - 1
	if i < 0 {
		return nil
	}

	name := names[i]
	return &optlevel.Named{Name: name, Level: opts[name]}

}

func (c *compilerPerturber) chooseMOpt(opts stringhelp.Set) string {
	// Don't bother trying to select an mopt if there aren't any
	// TODO(@MattWindsor91): should this be an error?
	nopts := len(opts)
	if nopts == 0 {
		return ""
	}
	// 'don't choose an mopt' - the empty string, may or may not be a valid choice, so we don't factor it in here.
	optsl := opts.Slice()
	i := c.rng.Intn(nopts)
	return optsl[i]
}
