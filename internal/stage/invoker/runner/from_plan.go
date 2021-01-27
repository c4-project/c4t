// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"github.com/c4-project/c4t/internal/copier"
	"github.com/c4-project/c4t/internal/plan"
	"github.com/c4-project/c4t/internal/remote"
)

// FromPlanFactory is a runner factory that instantiates either a SSH or local runner depending on the machine
// configuration inside the first plan passed to it.
//
// This is useful for single-shot invocation over a plan, where there is no benefit to setting up a connection based
// on central machine/SSH configuration.  In the director, the invoker will set up the machine configuration in
// advance, and there is no need to consult the plan.
type FromPlanFactory struct {
	// The global remoting config used for any remote connections initiated by this factory.
	Config *remote.Config

	cached Factory
}

// MakeRunner makes a runner using the machine configuration in pl.
func (p *FromPlanFactory) MakeRunner(ldir string, pl *plan.Plan, obs ...copier.Observer) (Runner, error) {
	var err error
	if p.cached == nil {
		if p.cached, err = p.makeFactory(pl); err != nil {
			return nil, err
		}
	}
	return p.cached.MakeRunner(ldir, pl, obs...)
}

func (p *FromPlanFactory) makeFactory(pl *plan.Plan) (Factory, error) {
	return FactoryFromRemoteConfig(p.Config, pl.Machine.SSH)
}

// Close closes the runner factory, if it was ever instantiated.
func (p *FromPlanFactory) Close() error {
	if p.cached == nil {
		return nil
	}
	return p.cached.Close()
}
