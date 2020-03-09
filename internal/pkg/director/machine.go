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

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/planner"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"github.com/MattWindsor91/act-tester/internal/pkg/config"
	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
)

// Machine contains the state necessary to run a single machine loop of a director.
type Machine struct {
	// Config contains the machine config for this machine.
	Config config.Machine

	// ID is the ID for this machine.
	ID model.ID

	// InFiles is the list of files to use as the base corpus for this machine loop.
	InFiles []string

	// Env contains the parts of the director's config that tell it how to do various environmental tasks.
	Env *Env

	// Logger points to a logger for this machine's loop.
	Logger *log.Logger

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

	if m.Config.SSH != nil {
		return errors.New("TODO: SSH support not yet available")
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

	if p, err = m.plan(ctx); err != nil {
		return fmt.Errorf("while planning: %w", err)
	}
	if p, err = m.fuzz(ctx, p); err != nil {
		return fmt.Errorf("while fuzzing: %w", err)
	}
	if p, err = m.lift(ctx, p); err != nil {
		return fmt.Errorf("while lifting: %w", err)
	}
	if _, err = m.mach(ctx, p); err != nil {
		return fmt.Errorf("while performing machine-specific actions: %w", err)
	}

	return nil
}

func (m *Machine) plan(ctx context.Context) (*plan.Plan, error) {
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
		Observer:  ux.NewPbObserver(m.Logger),
		InFiles:   m.InFiles,
		MachineID: m.ID,
	}
	return &p, nil
}

func (m *Machine) fuzz(_ context.Context, _ *plan.Plan) (*plan.Plan, error) {
	return nil, nil
}

func (m *Machine) lift(_ context.Context, _ *plan.Plan) (*plan.Plan, error) {
	return nil, nil
}

func (m *Machine) mach(_ context.Context, _ *plan.Plan) (*plan.Plan, error) {
	return nil, nil
}
