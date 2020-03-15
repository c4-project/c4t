// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/planner"

	"github.com/MattWindsor91/act-tester/internal/pkg/config"
	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
)

// Machine contains the state necessary to run a single machine loop of a director.
type Machine struct {
	// MachConfig contains the machine config for this machine.
	MachConfig config.Machine

	// FuzzConfig is the configuration for this machine's fuzzer.
	FuzzConfig *fuzzer.Config

	// ID is the ID for this machine.
	ID model.ID

	// InFiles is the list of files to use as the base corpus for this machine loop.
	InFiles []string

	// Env contains the parts of the director's config that tell it how to do various environmental tasks.
	Env *Env

	// Logger points to a logger for this machine's loop.
	Logger *log.Logger

	// Observer is this machine's builder observer.
	Observer builder.Observer

	// Paths contains the scratch pathset for this machine.
	Paths *MachinePathset
}

// Run runs this machine's testing loop.
func (m *Machine) Run(ctx context.Context) error {
	m.Logger = iohelp.EnsureLog(m.Logger)
	if err := m.check(); err != nil {
		return err
	}

	m.Logger.Print("preparing scratch directories")
	if err := m.Paths.Prepare(); err != nil {
		return err
	}

	m.Logger.Print("starting loop")
	return m.mainLoop(ctx)
}

// check makes sure this machine has a valid configuration before starting loops.
func (m *Machine) check() error {
	if m.Paths == nil {
		return fmt.Errorf("%w: paths for machine %s", iohelp.ErrPathsetNil, m.ID.String())
	}

	if m.Env == nil {
		return errors.New("no environment configuration")
	}

	if m.FuzzConfig == nil {
		return errors.New("no fuzzer config")
	}

	if m.MachConfig.SSH != nil {
		return errors.New("TODO: SSH support not yet available")
	}

	if m.Observer == nil {
		return errors.New("observer nil")
	}

	return nil
}

// mainLoop performs the main testing loop for one machine.
func (m *Machine) mainLoop(ctx context.Context) error {
	for {
		if err := m.pass(ctx); err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
	}
}

// pass performs one iteration of the main testing loop for one machine.
func (m *Machine) pass(ctx context.Context) error {
	var (
		p   *plan.Plan
		err error
	)

	steps := []struct {
		name string
		f    func(context.Context, *plan.Plan) (*plan.Plan, error)
	}{
		{name: "init", f: m.plan},
		{name: "fuzz", f: m.fuzz},
		//{name: "lift", f: m.lift},
	}

	for _, s := range steps {
		if p, err = s.f(ctx, p); err != nil {
			return fmt.Errorf("in %s stage: %w", s.name, err)
		}
		if err = m.dump(s.name, p); err != nil {
			return fmt.Errorf("when dumping after %s stage: %w", s.name, err)
		}
	}

	return nil
}

// plan creates a new plan using ctx.
// It ignores the incoming plan; in practice, this will be a nil pointer,
// and exists only to bring this method in line with the signatures of the other pass methods.
func (m *Machine) plan(ctx context.Context, _ *plan.Plan) (*plan.Plan, error) {
	p, err := m.makePlanner()
	if err != nil {
		return nil, err
	}
	return p.Plan(ctx)
}

func (m *Machine) makePlanner() (*planner.Planner, error) {
	p := planner.Planner{
		Source:    m.Env.Planner,
		Logger:    m.Logger,
		Observer:  m.Observer,
		InFiles:   m.InFiles,
		MachineID: m.ID,
	}
	return &p, nil
}

func (m *Machine) fuzz(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	f, err := fuzzer.New(m.FuzzConfig, p)
	if err != nil {
		return nil, err
	}
	return f.Fuzz(ctx)
}

// dump dumps a plan p to its expected plan file given the stage name name.
func (m *Machine) dump(name string, p *plan.Plan) error {
	fname := m.Paths.PlanForStage(name)
	f, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("while opening plan file for %s: %w", name, err)
	}
	if err := p.Dump(f); err != nil {
		_ = f.Close()
		return fmt.Errorf("while writing plan file for %s: %w", name, err)
	}
	return f.Close()
}

/*
func (m *Machine) lift(_ context.Context, _ *plan.Plan) (*plan.Plan, error) {
	return nil, nil
}

func (m *Machine) mach(_ context.Context, _ *plan.Plan) (*plan.Plan, error) {
	return nil, nil
}
*/
