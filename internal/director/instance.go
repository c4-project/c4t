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
	"time"

	"github.com/MattWindsor91/act-tester/internal/stage/invoker/runner"

	observer2 "github.com/MattWindsor91/act-tester/internal/stage/mach/observer"

	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/stage/perturber"

	"github.com/MattWindsor91/act-tester/internal/helper/errhelp"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/copier"
	"github.com/MattWindsor91/act-tester/internal/machine"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"

	"github.com/MattWindsor91/act-tester/internal/model/run"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/director/observer"

	"github.com/MattWindsor91/act-tester/internal/director/pathset"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/remote"

	"github.com/MattWindsor91/act-tester/internal/stage/invoker"

	"github.com/MattWindsor91/act-tester/internal/stage/lifter"

	"github.com/MattWindsor91/act-tester/internal/stage/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/plan"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

// The maximum permitted number of times a loop can error out consecutively before the tester fails.
const maxConsecutiveErrors = 10

// Instance contains the state necessary to run a single machine loop of a director.
type Instance struct {
	// MachConfig contains the machine config for this machine.
	MachConfig machine.Config
	// SSHConfig contains top-level SSH configuration.
	SSHConfig *remote.Config
	// StageConfig is the configuration for this instance's stages.
	StageConfig *StageConfig

	// ID is the ID for this machine.
	ID id.ID

	// InitialPlan is the plan that is perturbed to form the plan for each test cycle.
	InitialPlan plan.Plan

	// Env contains the parts of the director's config that tell it how to do various environmental tasks.
	Env *Env

	// Logger points to a logger for this machine's loop.
	Logger *log.Logger

	// Observers is this machine's observer set.
	Observers []observer.Instance

	// SavedPaths contains the save pathset for this machine.
	SavedPaths *saver.Pathset
	// ScratchPaths contains the scratch pathset for this machine.
	ScratchPaths *pathset.Scratch

	// Quantities contains the quantity set for this machine.
	Quantities quantity.MachineSet
}

// Run runs this machine's testing loop.
func (i *Instance) Run(ctx context.Context) error {
	i.Logger = iohelp.EnsureLog(i.Logger)
	if err := i.check(); err != nil {
		return err
	}

	i.Logger.Println("preparing scratch directories")
	if err := i.ScratchPaths.Prepare(); err != nil {
		return err
	}

	i.Logger.Println("creating stage configurations")
	sc, err := i.makeStageConfig()
	if err != nil {
		return err
	}
	i.Logger.Println("checking stage configurations")
	if err := sc.Check(); err != nil {
		return err
	}

	i.Logger.Println("starting loop")
	err = i.mainLoop(ctx, sc)
	i.Logger.Println("cleaning up")
	cerr := i.cleanUp()
	return errhelp.FirstError(err, cerr)
}

// cleanUp closes things that should be gracefully closed after an instance terminates.
func (i *Instance) cleanUp() error {
	if i.StageConfig != nil && i.StageConfig.Invoke != nil {
		return i.StageConfig.Invoke.Close()
	}
	return nil
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

	return nil
}

// mainLoop performs the main testing loop for one machine.
func (i *Instance) mainLoop(ctx context.Context, sc *StageConfig) error {
	var (
		iter    uint64
		nErrors uint
	)
	for {
		if err := i.iterate(ctx, iter, sc); err != nil {
			// This serves to stop the tester if we get stuck in a rapid failure loop on a particular machine.
			// TODO(@MattWindsor91): ideally this should be timing the gap between errors, so that we stop if there
			// are too many errors happening too quickly.
			nErrors++
			if maxConsecutiveErrors < nErrors {
				return fmt.Errorf("too many consecutive errors; last error was: %w", err)
			}
			i.Logger.Println("ERROR:", err)
		} else {
			nErrors = 0
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		iter++
	}
}

// iterate performs one iteration of the main testing loop (number iter) for one machine.
func (i *Instance) iterate(ctx context.Context, iter uint64, sc *StageConfig) error {
	var err error

	pcopy := i.InitialPlan
	p := &pcopy

	r := run.Run{
		MachineID: i.ID,
		Iter:      iter,
		Start:     time.Now(),
	}
	observer.OnIteration(r, i.Observers...)

	for _, s := range Stages {
		p, err = i.runStage(ctx, sc, s, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Instance) runStage(ctx context.Context, sc *StageConfig, s stageRunner, p *plan.Plan) (*plan.Plan, error) {
	var err error
	if p, err = s.Run(sc, ctx, p); err != nil {
		return nil, fmt.Errorf("in %s stage: %w", s.Stage, err)
	}
	if err = i.dump(s.Stage, p); err != nil {
		return nil, fmt.Errorf("when dumping after %s stage: %w", s.Stage, err)
	}
	return p, nil
}

func (i *Instance) makeStageConfig() (*StageConfig, error) {
	bobs := observer.LowerToBuilder(i.Observers)
	cobs := observer.LowerToCopy(i.Observers)

	var (
		err error
		sc  StageConfig
	)

	if sc.Perturb, err = i.makePerturber(observer.LowerToPerturber(i.Observers)); err != nil {
		return nil, fmt.Errorf("when making planner: %w", err)
	}
	if sc.Fuzz, err = i.makeFuzzer(bobs); err != nil {
		return nil, fmt.Errorf("when making fuzzer config: %w", err)
	}
	if sc.Lift, err = i.makeLifter(bobs); err != nil {
		return nil, fmt.Errorf("when making lifter config: %w", err)
	}
	if sc.Invoke, err = i.makeInvoker(cobs, observer.LowerToMach(i.Observers)); err != nil {
		return nil, fmt.Errorf("when making machine invoker: %w", err)
	}
	if sc.Analyser, err = i.makeAnalyser(observer.LowerToAnalyser(i.Observers), observer.LowerToSaver(i.Observers)); err != nil {
		return nil, fmt.Errorf("when making analysis: %w", err)
	}
	return &sc, nil
}

func (i *Instance) makeAnalyser(aobs []analyser.Observer, sobs []saver.Observer) (*analyser.Analyser, error) {
	return analyser.New(
		analyser.ObserveWith(aobs...),
		analyser.ObserveSaveWith(sobs...),
		analyser.ParWorkers(10), // TODO(@MattWindsor91): get this from somewhere
		analyser.SaveToPathset(i.SavedPaths),
	)
}

func (i *Instance) makePerturber(obs []perturber.Observer) (*perturber.Perturber, error) {
	return perturber.New(
		i.Env.CInspector,
		perturber.ObserveWith(obs...),
		perturber.OverrideQuantities(i.Quantities.Perturb),
	)
}

func (i *Instance) makeFuzzer(obs []builder.Observer) (*fuzzer.Fuzzer, error) {
	return fuzzer.New(
		i.Env.Fuzzer,
		fuzzer.NewPathset(i.ScratchPaths.DirFuzz),
		fuzzer.ObserveWith(obs...),
		fuzzer.LogWith(i.Logger),
		fuzzer.OverrideQuantities(i.Quantities.Fuzz),
	)
}

func (i *Instance) makeLifter(obs []builder.Observer) (*lifter.Lifter, error) {
	return lifter.New(
		i.Env.Lifter,
		lifter.NewPathset(i.ScratchPaths.DirLift),
		lifter.LogTo(i.Logger),
		lifter.ObserveWith(obs...),
	)
}

func (i *Instance) makeInvoker(cobs []copier.Observer, mobs []observer2.Observer) (*invoker.Invoker, error) {
	// Unlike the single-shot, we don't late-bind the factory using the plan.  This is because we've already
	// got the machine configuration without it.
	f, err := runner.FactoryFromRemoteConfig(i.SSHConfig, i.MachConfig.SSH)
	if err != nil {
		return nil, err
	}
	return invoker.New(i.ScratchPaths.DirRun,
		f,
		invoker.ObserveCopiesWith(cobs...),
		invoker.ObserveMachWith(mobs...),
		// As above, there is no loading of quantities using the plan, as we already know which machine the plan is
		// targeting without consulting the plan.
		invoker.OverrideBaseQuantities(i.Quantities.Mach),
	)
}

// dump dumps a plan p to its expected plan file given the stage s.
func (i *Instance) dump(s stage.Stage, p *plan.Plan) error {
	// TODO(@MattWindsor91): is this really necessary anymore?
	return p.WriteFile(i.ScratchPaths.PlanForStage(s), plan.WriteHuman)
}
