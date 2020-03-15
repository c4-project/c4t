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

	"github.com/MattWindsor91/act-tester/internal/pkg/lifter"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/planner"

	"github.com/MattWindsor91/act-tester/internal/pkg/config"
	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
)

// Instance contains the state necessary to run a single machine loop of a director.
type Instance struct {
	// MachConfig contains the machine config for this machine.
	MachConfig config.Machine

	// StageConfig is the configuration for this instance's stages.
	StageConfig *StageConfig

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

	// Quantities contains the quantity set for this machine.
	Quantities config.QuantitySet
}

// Run runs this machine's testing loop.
func (m *Instance) Run(ctx context.Context) error {
	m.Logger = iohelp.EnsureLog(m.Logger)
	if err := m.check(); err != nil {
		return err
	}

	m.Logger.Print("preparing scratch directories")
	if err := m.Paths.Prepare(); err != nil {
		return err
	}

	m.Logger.Print("creating stage configurations")
	sc, err := m.makeStageConfig()
	if err != nil {
		return err
	}
	m.Logger.Print("checking stage configurations")
	if err := sc.Check(); err != nil {
		return err
	}

	m.Logger.Print("starting loop")
	return m.mainLoop(ctx, sc)
}

// check makes sure this machine has a valid configuration before starting loops.
func (m *Instance) check() error {
	if m.Paths == nil {
		return fmt.Errorf("%w: paths for machine %s", iohelp.ErrPathsetNil, m.ID.String())
	}

	if m.Env == nil {
		return errors.New("no environment configuration")
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
func (m *Instance) mainLoop(ctx context.Context, sc *StageConfig) error {
	for {
		if err := m.pass(ctx, sc); err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
	}
}

// pass performs one iteration of the main testing loop for one machine.
func (m *Instance) pass(ctx context.Context, sc *StageConfig) error {
	var (
		p   *plan.Plan
		err error
	)

	for _, s := range Stages {
		if p, err = s.Run(sc, ctx, p); err != nil {
			return fmt.Errorf("in %s stage: %w", s.Name, err)
		}
		if err = m.dump(s.Name, p); err != nil {
			return fmt.Errorf("when dumping after %s stage: %w", s.Name, err)
		}
	}

	return nil
}

func (m *Instance) makeStageConfig() (*StageConfig, error) {
	p, err := m.makePlanner()
	if err != nil {
		return nil, fmt.Errorf("when making planner: %w", err)
	}
	f, err := m.makeFuzzerConfig()
	if err != nil {
		return nil, fmt.Errorf("when making fuzzer config: %w", err)
	}
	l, err := m.makeLifterConfig()
	if err != nil {
		return nil, fmt.Errorf("when making lifter config: %w", err)
	}
	c := &LocalMach{Dir: m.Paths.DirRun}
	sc := StageConfig{
		InFiles: m.InFiles,
		Plan:    p,
		Fuzz:    f,
		Lift:    l,
		Mach:    c,
	}
	return &sc, nil
}

func (m *Instance) makePlanner() (*planner.Planner, error) {
	p := planner.Planner{
		Source:    m.Env.Planner,
		Logger:    m.Logger,
		Observer:  m.Observer,
		MachineID: m.ID,
	}
	return &p, nil
}

func (m *Instance) makeFuzzerConfig() (*fuzzer.Config, error) {
	fz := m.Env.Fuzzer
	if fz == nil {
		return nil, errors.New("no single fuzzer provided")
	}

	fc := fuzzer.Config{
		Driver:     fz,
		Logger:     m.Logger,
		Observer:   m.Observer,
		Paths:      fuzzer.NewPathset(m.Paths.DirFuzz),
		Quantities: m.Quantities.Fuzz,
	}

	return &fc, nil
}

func (m *Instance) makeLifterConfig() (*lifter.Config, error) {
	hm := m.Env.Lifter
	if hm == nil {
		return nil, errors.New("no single fuzzer provided")
	}

	lc := lifter.Config{
		Maker:    hm,
		Logger:   m.Logger,
		Observer: m.Observer,
		Paths:    lifter.NewPathset(m.Paths.DirLift),
	}

	return &lc, nil
}

// dump dumps a plan p to its expected plan file given the stage name name.
func (m *Instance) dump(name string, p *plan.Plan) error {
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
