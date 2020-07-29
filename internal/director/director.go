// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package director contains the top-level ACT test director, which manages a full testing campaign.
package director

import (
	"context"
	"fmt"
	"log"

	"github.com/MattWindsor91/act-tester/internal/helper/errhelp"

	"github.com/MattWindsor91/act-tester/internal/model/machine"

	"github.com/MattWindsor91/act-tester/internal/director/pathset"
	"github.com/MattWindsor91/act-tester/internal/remote"

	"github.com/MattWindsor91/act-tester/internal/director/observer"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/config"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

// Director contains the main state and configuration for the test director.
type Director struct {
	// paths provides path resolving functionality for the director.
	paths *pathset.Pathset
	// machines contains the machines that will be used in the test.
	machines machine.ConfigMap
	// observers contains multi-machine observers for the director.
	observers []observer.Observer
	// env groups together the bits of configuration that pertain to dealing with the environment.
	env Env
	// ssh, if present, provides configuration for the director's remote invocation.
	ssh *remote.Config
	// quantities contains various tunable quantities for the director's stages.
	quantities config.QuantitySet
	// files is the input file set.
	files []string
	// l is the logger for the director.
	l *log.Logger
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
	return &d, d.tidyAfterOptions()
}

func (d *Director) tidyAfterOptions() error {
	if len(d.machines) == 0 {
		return ErrNoMachines
	}
	if d.paths == nil {
		return iohelp.ErrPathsetNil
	}
	d.l = iohelp.EnsureLog(d.l)
	return nil
}

// liftInitError lifts err to mention that it occurred during initialisation of a director.
func liftInitError(err error) error {
	return fmt.Errorf("while initialising director: %w", err)
}

// Direct runs the director d, closing all of its observers on termination.
func (d *Director) Direct(ctx context.Context) error {
	err := d.directInner(ctx)
	cerr := observer.CloseAll(d.observers...)
	return errhelp.FirstError(err, cerr)
}

func (d *Director) directInner(ctx context.Context) error {
	if err := d.prepare(); err != nil {
		return err
	}

	ms, err := d.makeMachines()
	if err != nil {
		return err
	}

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return d.runLoops(cctx, cancel, ms)
}

func (d *Director) runLoops(cctx context.Context, cancel func(), ms []*Instance) error {
	eg, ectx := errgroup.WithContext(cctx)
	for _, m := range ms {
		m := m
		eg.Go(func() error { return m.Run(ectx) })
	}
	for _, o := range d.observers {
		o := o
		eg.Go(func() error { return o.Run(ectx, cancel) })
	}
	return eg.Wait()
}

func (d *Director) makeMachines() ([]*Instance, error) {
	ms := make([]*Instance, len(d.machines))
	var (
		i   int
		err error
	)
	for midstr, c := range d.machines {
		if ms[i], err = d.makeMachine(midstr, c); err != nil {
			return nil, err
		}
		i++
	}
	return ms, nil
}

func (d *Director) prepare() error {
	d.quantities.Log(d.l)

	d.l.Println("making directories")
	if err := d.paths.Prepare(); err != nil {
		return err
	}

	return d.machines.ObserveOn(observer.LowerToMachine(d.observers)...)
}

func (d *Director) makeMachine(midstr string, c machine.Config) (*Instance, error) {
	l := log.New(d.l.Writer(), logPrefix(midstr), 0)
	mid, err := id.TryFromString(midstr)
	if err != nil {
		return nil, err
	}
	obs, err := d.instanceObservers(mid)
	if err != nil {
		return nil, err
	}
	sps := d.paths.MachineScratch(mid)
	vps := d.paths.MachineSaved(mid)
	m := Instance{
		MachConfig:   c,
		SSHConfig:    d.ssh,
		Env:          &d.env,
		ID:           mid,
		InFiles:      d.files,
		Observers:    obs,
		ScratchPaths: sps,
		SavedPaths:   vps,
		Logger:       l,
		Quantities:   d.quantities,
	}
	return &m, nil
}

func (d *Director) instanceObservers(mid id.ID) ([]observer.Instance, error) {
	var err error
	ios := make([]observer.Instance, len(d.observers))
	for i, o := range d.observers {
		if ios[i], err = o.Instance(mid); err != nil {
			return nil, err
		}
	}
	return ios, nil
}

func logPrefix(midstr string) string {
	return midstr + ": "
}
