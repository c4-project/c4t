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

	copy2 "github.com/MattWindsor91/act-tester/internal/copier"

	"github.com/MattWindsor91/act-tester/internal/model/machine"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse"

	aobserver "github.com/MattWindsor91/act-tester/internal/controller/analyse/observer"
	"github.com/MattWindsor91/act-tester/internal/controller/analyse/saver"

	"github.com/MattWindsor91/act-tester/internal/controller/mach"
	"github.com/MattWindsor91/act-tester/internal/view/stdflag"

	"github.com/MattWindsor91/act-tester/internal/model/run"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/director/observer"

	"github.com/MattWindsor91/act-tester/internal/director/pathset"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/remote"

	"github.com/MattWindsor91/act-tester/internal/controller/rmach"

	"github.com/MattWindsor91/act-tester/internal/controller/lifter"

	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/controller/planner"
	"github.com/MattWindsor91/act-tester/internal/model/plan"

	"github.com/MattWindsor91/act-tester/internal/config"
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

	// InFiles is the list of files to use as the base corpus for this machine loop.
	InFiles []string

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
	Quantities config.QuantitySet
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
	return iohelp.FirstError(err, cerr)
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
	var (
		p   *plan.Plan
		err error
	)

	r := run.Run{
		MachineID: i.ID,
		Iter:      iter,
		Start:     time.Now(),
	}
	observer.OnIteration(r, i.Observers...)

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
	aobs := observer.LowerToAnalyse(i.Observers)
	bobs := observer.LowerToBuilder(i.Observers)
	cobs := observer.LowerToCopy(i.Observers)

	var (
		err error
		sc  StageConfig
	)

	if sc.Plan, err = i.makePlanner(observer.LowerToPlanner(i.Observers)); err != nil {
		return nil, fmt.Errorf("when making planner: %w", err)
	}
	if sc.Fuzz, err = i.makeFuzzerConfig(bobs); err != nil {
		return nil, fmt.Errorf("when making fuzzer config: %w", err)
	}
	if sc.Lift, err = i.makeLifterConfig(bobs); err != nil {
		return nil, fmt.Errorf("when making lifter config: %w", err)
	}
	if sc.Invoke, err = i.makeInvoker(cobs, bobs); err != nil {
		return nil, fmt.Errorf("when making machine invoker: %w", err)
	}
	sc.Analyse = i.makeAnalyseConfig(aobs)
	return &sc, nil
}

func (i *Instance) makeAnalyseConfig(aobs []aobserver.Observer) *analyse.Config {
	return &analyse.Config{
		Observers:  aobs,
		NWorkers:   10, // TODO(@MattWindsor91): get this from somewhere
		SavedPaths: i.SavedPaths,
	}
}

func (i *Instance) makePlanner(obs []planner.Observer) (*planner.Planner, error) {
	// TODO(@MattWindsor91): move planner config outside of instance
	c := planner.Config{
		Source:     i.Env.Planner,
		Logger:     i.Logger,
		Observers:  planner.NewObserverSet(obs...),
		Quantities: i.Quantities.Plan,
	}
	return planner.New(&c, i.machineForPlan(), i.InFiles, plan.UseDateSeed)
}

// machineForPlan massages this instance's machine config into a form with which the planner is comfortable.
func (i *Instance) machineForPlan() machine.Named {
	return machine.Named{
		ID:      i.ID,
		Machine: i.MachConfig.Machine,
	}
}

func (i *Instance) makeFuzzerConfig(obs []builder.Observer) (*fuzzer.Config, error) {
	fz := i.Env.Fuzzer
	if fz == nil {
		return nil, errors.New("no single fuzzer provided")
	}

	fc := fuzzer.Config{
		Driver:     fz,
		Logger:     i.Logger,
		Observers:  obs,
		Paths:      fuzzer.NewPathset(i.ScratchPaths.DirFuzz),
		Quantities: i.Quantities.Fuzz,
	}

	return &fc, nil
}

func (i *Instance) makeLifterConfig(obs []builder.Observer) (*lifter.Config, error) {
	hm := i.Env.Lifter
	if hm == nil {
		return nil, errors.New("no single lifter provided")
	}

	lc := lifter.Config{
		Driver:    hm,
		Logger:    i.Logger,
		Observers: obs,
		Paths:     lifter.NewPathset(i.ScratchPaths.DirLift),
	}

	return &lc, nil
}

func (i *Instance) makeInvoker(cobs []copy2.Observer, bobs []builder.Observer) (*rmach.Invoker, error) {
	return rmach.New(i.ScratchPaths.DirRun,
		stdflag.MachInvoker{
			// TODO(@MattWindsor91): this is a bit messy.
			Config: &mach.UserConfig{
				OutDir:     i.ScratchPaths.DirRun,
				Quantities: i.Quantities.Mach,
			},
		},
		rmach.ObserveCopiesWith(cobs...),
		rmach.ObserveCorpusWith(bobs...),
		rmach.UseSSH(i.SSHConfig, i.MachConfig.SSH),
	)
}

// dump dumps a plan p to its expected plan file given the stage name name.
func (i *Instance) dump(name string, p *plan.Plan) error {
	return p.WriteFile(i.ScratchPaths.PlanForStage(name))
}
