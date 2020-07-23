// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"github.com/MattWindsor91/act-tester/internal/copier"
	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/remote"
)

// FromPlanFactory is a runner factory that instantiates either a SSH or local runner depending on the machine
// configuration inside the first plan passed to it.
//
// This is useful for single-shot invocation over a plan, where there is no benefit to setting up a connection based
// on central machine/SSH configuration.
type FromPlanFactory struct {
	// The local root directory to use for invocation results.
	LocalRoot string
	// The global remoting config used for any remote connections initiated by this factory.
	Config *remote.Config

	cached Factory
}

// MakeRunner makes a runner using the machine configuration in pl.
func (p *FromPlanFactory) MakeRunner(pl *plan.Plan, obs ...copier.Observer) (Runner, error) {
	var err error
	if p.cached == nil {
		if p.cached, err = p.makeFactory(pl); err != nil {
			return nil, err
		}
	}
	return p.cached.MakeRunner(pl, obs...)
}

func (p *FromPlanFactory) makeFactory(pl *plan.Plan) (Factory, error) {
	if pl.Machine.SSH == nil {
		return LocalFactory(p.LocalRoot), nil
	}
	return NewRemoteFactory(p.LocalRoot, p.Config, pl.Machine.SSH)
}

// Close closes the runner factory, if it was ever instantiated.
func (p *FromPlanFactory) Close() error {
	if p.cached == nil {
		return nil
	}
	return p.cached.Close()
}
