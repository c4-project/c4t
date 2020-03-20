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
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/director/pathset"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/remote"

	"github.com/MattWindsor91/act-tester/internal/pkg/director/mach"

	"github.com/MattWindsor91/act-tester/internal/pkg/lifter"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/planner"

	"github.com/MattWindsor91/act-tester/internal/pkg/config"
	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
)

// The maximum permitted number of times a loop can error out consecutively before the tester fails.
const maxConsecutiveErrors = 10

// Instance contains the state necessary to run a single machine loop of a director.
type Instance struct {
	// MachConfig contains the machine config for this machine.
	MachConfig config.Machine
	// SSHConfig contains top-level SSH configuration.
	SSHConfig *remote.Config
	// StageConfig is the configuration for this instance's stages.
	StageConfig *StageConfig

	// ID is the ID for this machine.
	ID id.ID

	// InFiles is the list of files to use as the base corpus for this machine loop.
	InFiles []string

	// Env contains the parts of the director's config that tell it how to do various environmental tasks.
	Env *Env

	// Logger points to a logger for this machine's loop.
	Logger *log.Logger

	// Observer is this machine's observer.
	Observer MachineObserver

	// SavedPaths contains the save pathset for this machine.
	SavedPaths *pathset.Saved
	// ScratchPaths contains the scratch pathset for this machine.
	ScratchPaths *pathset.Scratch

	// Quantities contains the quantity set for this machine.
	Quantities config.QuantitySet
}

// Run runs this machine's testing loop.
func (i *Instance) Run(ctx context.Context) error {
	i.Logger = iohelp.EnsureLog(i.Logger)
	if err := i.check(); err != nil {
		return err
	}

	i.Logger.Print("preparing scratch directories")
	if err := i.ScratchPaths.Prepare(); err != nil {
		return err
	}

	i.Logger.Print("creating stage configurations")
	sc, err := i.makeStageConfig()
	if err != nil {
		return err
	}
	i.Logger.Print("checking stage configurations")
	if err := sc.Check(); err != nil {
		return err
	}

	i.Logger.Print("starting loop")
	return i.mainLoop(ctx, sc)
}

// check makes sure this machine has a valid configuration before starting loops.
func (i *Instance) check() error {
	if i.ScratchPaths == nil {
		return fmt.Errorf("%w: paths for machine %s", iohelp.ErrPathsetNil, i.ID.String())
	}

	if i.Env == nil {
		return errors.New("no environment configuration")
	}

	// TODO(@MattWindsor): check SSHConfig?

	if i.Observer == nil {
		return errors.New("observer nil")
	}

	return nil
}

// mainLoop performs the main testing loop for one machine.
func (i *Instance) mainLoop(ctx context.Context, sc *StageConfig) error {
	var (
		iter    uint64
		nErrors uint
	)
	for {
		if err := i.pass(ctx, iter, sc); err != nil {
			// This serves to stop the tester if we get stuck in a rapid failure loop on a particular machine.
			// TODO(@MattWindsor91): ideally this should be timing the gap between errors, so that we stop if there
			// are too many errors happening too quickly.
			nErrors++
			if maxConsecutiveErrors < nErrors {
				return fmt.Errorf("too many consecutive errors; last error was: %w", err)
			}
			i.Logger.Println("ERROR:", err)
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		iter++
		nErrors = 0
	}
}

// pass performs one iteration of the main testing loop (number iter) for one machine.
func (i *Instance) pass(ctx context.Context, iter uint64, sc *StageConfig) error {
	var (
		p   *plan.Plan
		err error
	)

	i.Observer.OnIteration(iter, time.Now())

	for _, s := range Stages {
		if p, err = s.Run(sc, ctx, p); err != nil {
			return fmt.Errorf("in %s stage: %w", s.Name, err)
		}
		if err = i.dump(s.Name, p); err != nil {
			return fmt.Errorf("when dumping after %s stage: %w", s.Name, err)
		}
	}

	return nil
}

func (i *Instance) makeStageConfig() (*StageConfig, error) {
	p, err := i.makePlanner()
	if err != nil {
		return nil, fmt.Errorf("when making planner: %w", err)
	}
	f, err := i.makeFuzzerConfig()
	if err != nil {
		return nil, fmt.Errorf("when making fuzzer config: %w", err)
	}
	l, err := i.makeLifterConfig()
	if err != nil {
		return nil, fmt.Errorf("when making lifter config: %w", err)
	}
	m, err := mach.New(i.Observer, i.ScratchPaths.DirRun, i.SSHConfig, i.MachConfig.SSH)
	if err != nil {
		return nil, fmt.Errorf("when making machine-exec config: %w", err)
	}
	sc := StageConfig{
		InFiles: i.InFiles,
		Plan:    p,
		Fuzz:    f,
		Lift:    l,
		Mach:    m,
		Save: &Save{
			Logger:   i.Logger,
			NWorkers: 10, // TODO(@MattWindsor91): get this from somewhere
			Paths:    i.SavedPaths,
		},
	}
	return &sc, nil
}

func (i *Instance) makePlanner() (*planner.Planner, error) {
	p := planner.Planner{
		Source:    i.Env.Planner,
		Logger:    i.Logger,
		Observer:  i.Observer,
		MachineID: i.ID,
	}
	return &p, nil
}

func (i *Instance) makeFuzzerConfig() (*fuzzer.Config, error) {
	fz := i.Env.Fuzzer
	if fz == nil {
		return nil, errors.New("no single fuzzer provided")
	}

	fc := fuzzer.Config{
		Driver:     fz,
		Logger:     i.Logger,
		Observer:   i.Observer,
		Paths:      fuzzer.NewPathset(i.ScratchPaths.DirFuzz),
		Quantities: i.Quantities.Fuzz,
	}

	return &fc, nil
}

func (i *Instance) makeLifterConfig() (*lifter.Config, error) {
	hm := i.Env.Lifter
	if hm == nil {
		return nil, errors.New("no single fuzzer provided")
	}

	lc := lifter.Config{
		Maker:    hm,
		Logger:   i.Logger,
		Observer: i.Observer,
		Paths:    lifter.NewPathset(i.ScratchPaths.DirLift),
	}

	return &lc, nil
}

// dump dumps a plan p to its expected plan file given the stage name name.
func (i *Instance) dump(name string, p *plan.Plan) error {
	fname := i.ScratchPaths.PlanForStage(name)
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
