// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package invoker handles invoking the 'mach' stage in a network-transparent manner with the act-tester-mach binary.
package invoker

import (
	"github.com/1set/gut/ystring"
	"github.com/MattWindsor91/act-tester/internal/copier"
	"github.com/MattWindsor91/act-tester/internal/quantity"
	"github.com/MattWindsor91/act-tester/internal/stage/invoker/runner"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer"
)

// Invoker runs the machine-runner, through SSH if needed.
//
// Much of the invoker's behaviour is injected through two sources: a runner factory (which can either pass in cached
// SSH-or-lack-therof configuration, or delegate to the incoming machine plan), and a set of quantity overrides (that
// can either be fully pre-cached, or also delegate in part to the plan).  This set-up is intended to let the
// single-shot binary for the invoker rely almost entirely on information coming to it through the plan, while also
// letting the director set up the configuration in advance and bypass plan inspection.
type Invoker struct {
	// ldir is the local directory to which machine node files are to be copied.
	ldir string
	// baseQuantities contains the base quantity set, to be overridden by any other quantities given.
	baseQuantities quantity.MachNodeSet
	// pqo is a hook to let baseQuantities be overridden by information from the plan in single-shot mode.
	pqo PlanQuantityOverrider
	// copyObservers is the set of observers listening for file copying.
	copyObservers []copier.Observer
	// machObservers is the set of observers listening for remote corpus manipulations.
	machObservers []observer.Observer
	// rfac governs how the invoker will run the machine node when given a plan to invoke.
	rfac runner.Factory
	// allowReinvoke permits re-invokation on plans that already have a reinvoke stage.
	allowReinvoke bool
}

// New constructs a new Invoker with local directory ldir, runner factory fac, and options o.
func New(ldir string, fac runner.Factory, o ...Option) (*Invoker, error) {
	if ystring.IsBlank(ldir) {
		return nil, ErrDirEmpty
	}

	invoker := Invoker{ldir: ldir, rfac: fac, pqo: NopPlanQuantityOverrider{}}
	if err := Options(o...)(&invoker); err != nil {
		return nil, err
	}
	return &invoker, nil
}
