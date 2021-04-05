// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package director contains the top-level C4 test director, which manages a full testing campaign.
package director

import (
	"context"
	"fmt"
	"time"

	fuzzer2 "github.com/c4-project/c4t/internal/model/service/fuzzer"

	"github.com/c4-project/c4t/internal/plan/analysis"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/stage/planner"

	"github.com/c4-project/c4t/internal/plan"

	"github.com/c4-project/c4t/internal/machine"

	"github.com/c4-project/c4t/internal/director/pathset"
	"github.com/c4-project/c4t/internal/remote"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/subject/corpus"

	"golang.org/x/sync/errgroup"

	"github.com/c4-project/c4t/internal/helper/iohelp"
)

// Director contains the main state and configuration for the test director.
type Director struct {
	// paths provides path resolving functionality for the director.
	paths *pathset.Pathset
	// machines contains the machines that will be used in the test.
	machines machine.ConfigMap
	// observers contains multi-machine observers for the director.
	observers []Observer
	// instances contains the instances governed by the director.
	instances []Instance
	// env groups together the bits of configuration that pertain to dealing with the environment.
	env Env
	// ssh, if present, provides configuration for the director's remote invocation.
	ssh *remote.Config
	// fcfg, if present, provides fuzzer configuration.
	fcfg *fuzzer2.Config
	// quantities contains various tunable quantities for the director's stages.
	quantities quantity.RootSet
	// files is the input file set.
	files []string
	// filters is the set of compiled filter sets to use in analysis.
	filters analysis.FilterSet
}

// New creates a new Director with driver set e, input paths files, machines ms, and options opt.
func New(e Env, ms machine.ConfigMap, files []string, opt ...Option) (*Director, error) {
	if len(files) == 0 {
		return nil, liftInitError(corpus.ErrNone)
	}
	if len(ms) == 0 {
		return nil, liftInitError(ErrNoMachines)
	}
	if err := e.Check(); err != nil {
		return nil, liftInitError(err)
	}
	d := Director{files: files, env: e, machines: ms}
	if err := Options(opt...)(&d); err != nil {
		return nil, liftInitError(err)
	}
	return &d, d.initAfterOptions()
}

func (d *Director) initAfterOptions() error {
	if len(d.machines) == 0 {
		return ErrNoMachines
	}
	if d.paths == nil {
		return iohelp.ErrPathsetNil
	}
	return d.initInstances()
}

// initInstances performs the initial set-up of instances (before allocation of plan resources to them).
func (d *Director) initInstances() error {
	// TODO(@MattWindsor91): eventually decouple machines from instances.

	d.instances = make([]Instance, len(d.machines))
	i := 0

	// This is a bit weird, but necessary at the moment to solve a race involving instance observers.
	OnPrepare(PrepareInstancesMessage(len(d.instances)), LowerToPrepare(d.observers)...)

	for mid, mc := range d.machines {
		if err := d.initInstance(i, mid, mc); err != nil {
			return err
		}
		i++
	}
	return nil
}

func (d *Director) initInstance(i int, mid id.ID, mc machine.Config) error {
	obs, err := d.instanceObservers(mid)
	if err != nil {
		return err
	}
	m := Machine{
		ID:         mid,
		Config:     mc,
		Pathset:    d.paths.Instance(mid),
		Quantities: d.machineQuantities(&mc),
	}
	d.instances[i] = Instance{
		Index:        i,
		SSHConfig:    d.ssh,
		Env:          d.env,
		Observers:    obs,
		Machine:      &m,
		Filters:      d.filters,
		FuzzerConfig: d.fcfg,
	}
	return nil
}

func (d *Director) machineQuantities(c *machine.Config) quantity.MachineSet {
	if c.Quantities == nil {
		return d.quantities.MachineSet
	}
	qs := d.quantities.MachineSet
	qs.Override(*c.Quantities)
	return qs
}

func (d *Director) instanceObservers(mid id.ID) ([]InstanceObserver, error) {
	var err error
	ios := make([]InstanceObserver, len(d.observers))
	for i, o := range d.observers {
		if ios[i], err = o.Instance(mid); err != nil {
			return nil, err
		}
	}
	return ios, nil
}

// liftInitError lifts err to mention that it occurred during initialisation of a director.
func liftInitError(err error) error {
	return fmt.Errorf("while initialising director: %w", err)
}

// Direct runs the director d.
func (d *Director) Direct(ctx context.Context) error {
	if err := d.prepare(ctx); err != nil {
		return err
	}

	pn, err := d.plan(ctx)
	if err != nil {
		return err
	}

	return d.runLoops(ctx, pn)
}

func (d *Director) plan(ctx context.Context) (plan.Map, error) {
	p, err := d.makePlanner()
	if err != nil {
		return nil, fmt.Errorf("when making planner: %w", err)
	}
	return p.Plan(ctx, d.machines, d.files...)
}

func (d *Director) makePlanner() (*planner.Planner, error) {
	return planner.New(
		d.env.Planner,
		planner.ObserveWith(LowerToPlanner(d.observers)...),
		planner.OverrideQuantities(d.quantities.Plan),
	)
}

func (d *Director) runLoops(ctx context.Context, plans plan.Map) error {
	// Doing this here so that the perceived experiment length is always at or slightly above the timeout interval.
	start := time.Now()
	cctx, cancel := d.quantities.GlobalTimeout.OnContext(ctx)
	defer cancel()

	d.announceTimes(cctx, start)

	eg, ectx := errgroup.WithContext(cctx)
	for _, m := range d.instances {
		m := m
		m.Machine.InitialPlan = plans[m.Machine.ID]
		eg.Go(func() error { return m.Run(ectx) })
	}
	return eg.Wait()
}

func (d *Director) announceTimes(ctx context.Context, start time.Time) {
	obs := LowerToPrepare(d.observers)
	OnPrepare(PrepareStartMessage(start), obs...)
	if dl, ok := ctx.Deadline(); ok {
		OnPrepare(PrepareTimeoutMessage(dl), obs...)
	}
}

func (d *Director) prepare(ctx context.Context) error {
	obs := LowerToPrepare(d.observers)

	if dl, ok := ctx.Deadline(); ok {
		OnPrepare(PrepareTimeoutMessage(dl), obs...)
	}

	OnPrepare(PrepareQuantitiesMessage(d.quantities), obs...)

	OnPrepare(PreparePathsMessage(*d.paths), obs...)
	if err := d.paths.Prepare(); err != nil {
		return err
	}

	return d.machines.ObserveOn(LowerToMachine(d.observers)...)
}
