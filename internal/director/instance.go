// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/c4-project/c4t/internal/mutation"

	fuzzer2 "github.com/c4-project/c4t/internal/model/service/fuzzer"

	"github.com/c4-project/c4t/internal/plan/analysis"

	"github.com/c4-project/c4t/internal/stage/invoker/runner"

	observer2 "github.com/c4-project/c4t/internal/stage/mach/observer"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/stage/perturber"

	"github.com/c4-project/c4t/internal/helper/errhelp"

	"github.com/c4-project/c4t/internal/copier"
	"github.com/c4-project/c4t/internal/machine"

	"github.com/c4-project/c4t/internal/stage/analyser"

	"github.com/c4-project/c4t/internal/stage/analyser/saver"

	"github.com/c4-project/c4t/internal/subject/corpus/builder"

	"github.com/c4-project/c4t/internal/director/pathset"

	"github.com/c4-project/c4t/internal/model/id"

	"github.com/c4-project/c4t/internal/remote"

	"github.com/c4-project/c4t/internal/stage/invoker"

	"github.com/c4-project/c4t/internal/stage/lifter"

	"github.com/c4-project/c4t/internal/stage/fuzzer"

	"github.com/c4-project/c4t/internal/plan"

	"github.com/c4-project/c4t/internal/helper/iohelp"
)

// Instance contains the state necessary to run a single loop of a director.
type Instance struct {
	// Index is the index of the instance in the director.
	Index int
	// Env contains the parts of the director's config that tell it how to do various environmental tasks.
	Env Env
	// Machine is the machine installed into the instance.
	Machine *Machine
	// Observers is this machine's observer set.
	Observers []InstanceObserver
	// SSHConfig contains top-level SSH configuration.
	SSHConfig *remote.Config
	// Filters contains the precompiled filter set for this instance.
	Filters analysis.FilterSet

	// CycleHooks contains a number of callbacks that are executed before beginning a cycle.
	CycleHooks []func(*Instance) error

	// TODO(@MattWindsor91): this configuration should ideally be per-machine, and then should be moved to Machine.

	// FuzzerConfig contains the fuzzer config for this instance.
	FuzzerConfig *fuzzer2.Configuration

	// Mutant is the current mutant that is in use on this instance, if any.
	mutant mutation.Mutant

	// timeoutCh stores the current error cooldown channel, if any.
	// This is refreshed whenever an error occurs.
	timeoutCh <-chan time.Time

	// cycleCh stores the current cycle result channel, if any.
	// This is refreshed whenever a new cycle is launched.
	cycleCh <-chan cycleResult
}

// Machine contains the state for a particular machine attached to an instance.
type Machine struct {
	// ID is the ID for this machine.
	ID id.ID

	// InitialPlan is the plan that is perturbed to form the plan for each test cycle.
	InitialPlan plan.Plan

	// Pathset contains the pathset for this instance.
	Pathset *pathset.Instance

	// Quantities contains the quantity set for this machine.
	Quantities quantity.MachineSet

	// Config contains the machine config for this machine.
	Config machine.Config

	// cycle is the number of the current cycle for the machine.
	// This is held separately from the instance as an instance may (eventually) run cycles for multiple machines.
	cycle uint64

	// stageConfig is the configuration for this instance's stages.
	stageConfig *StageConfig
}

// Run runs this instance's testing loop.
func (i *Instance) Run(ctx context.Context) error {
	if err := i.check(); err != nil {
		return err
	}

	if err := i.Machine.Pathset.Scratch.Prepare(); err != nil {
		return err
	}

	var err error
	if i.Machine.stageConfig, err = i.makeStageConfig(); err != nil {
		return err
	}
	if err := i.Machine.stageConfig.Check(); err != nil {
		return err
	}

	err = i.mainLoop(ctx)

	OnInstanceClose(i.Observers...)
	cerr := i.cleanUp()

	return errhelp.FirstError(err, cerr)
}

// cleanUp closes things that should be gracefully closed after an instance terminates.
func (i *Instance) cleanUp() error {
	if i.Machine == nil {
		return nil
	}
	return i.Machine.cleanUp()
}

func (m *Machine) cleanUp() error {
	if m.stageConfig != nil && m.stageConfig.Invoke != nil {
		return m.stageConfig.Invoke.Close()
	}
	return nil
}

// check makes sure this instance has a valid configuration before starting loops.
func (i *Instance) check() error {
	if i.Machine == nil {
		return errors.New("machine nil")
	}
	if err := i.Machine.check(); err != nil {
		return err
	}

	// TODO(@MattWindsor): check SSHConfig?

	return i.Env.Check()
}

func (m *Machine) check() error {
	if m.Pathset == nil {
		return fmt.Errorf("%w: paths for machine %s", iohelp.ErrPathsetNil, m.ID.String())
	}
	return nil
}

type cycleResult struct {
	cycle Cycle
	err   error
}

// mainLoop performs the main testing loop for one machine.
func (i *Instance) mainLoop(ctx context.Context) error {
	i.launch(ctx)
	for {
		select {
		case <-ctx.Done():
			i.drainCycleCh()
			return ctx.Err()
		case res := <-i.cycleCh:
			i.handleCycleEnd(ctx, res)
		case <-i.timeoutCh:
			i.launch(ctx)
		}
	}
}

func (i *Instance) handleCycleEnd(ctx context.Context, res cycleResult) {
	i.cycleCh = nil
	if res.err != nil {
		i.handleError(res.err, res)
		return
	}
	// Don't clean up scratch after a failing iteration; we might need the information in the scratch
	if err := i.cleanUpCycle(); err != nil {
		i.handleError(err, res)
		return
	}
	// Only re-launch if we actually managed to complete the cycle without any errors; otherwise, wait on i.timeoutCh
	i.launch(ctx)
}

func (i *Instance) drainCycleCh() {
	if i.cycleCh == nil {
		return
	}
	for range i.cycleCh {
	}
}

func (i *Instance) handleError(err error, res cycleResult) {
	OnCycle(CycleErrorMessage(res.cycle, err), i.Observers...)
	i.timeoutCh = time.After(5 * time.Second)
	// TODO(@MattWindsor91): exponential backoff timeout
}

// launch launches one iteration of the main testing loop for one machine.
func (i *Instance) launch(ctx context.Context) {
	i.timeoutCh = nil

	c := i.makeCycleInstance()
	OnCycle(CycleStartMessage(c.cycle), i.Observers...)

	i.Machine.cycle++

	ch := make(chan cycleResult)
	go func() {
		err := c.run(ctx)
		select {
		case <-ctx.Done():
		case ch <- cycleResult{cycle: c.cycle, err: err}:
		}
		close(ch)
	}()

	i.cycleCh = ch
}

func (i *Instance) makeCycleInstance() cycleInstance {
	return cycleInstance{
		cycle: Cycle{
			Instance:  i.Index,
			MachineID: i.Machine.ID,
			Iter:      i.Machine.cycle,
			Start:     time.Now(),
		},
		p:  i.plan(),
		sc: i.Machine.stageConfig,
	}
}

func (i *Instance) plan() *plan.Plan {
	// Important to _copy_ the plan
	pcopy := i.Machine.InitialPlan
	if pcopy.IsMutationTest() {
		pcopy.Mutation.Selection = i.mutant
	}
	return &pcopy
}

func (i *Instance) makeStageConfig() (*StageConfig, error) {
	bobs := LowerToBuilder(i.Observers)
	cobs := LowerToCopy(i.Observers)

	var (
		err error
		sc  StageConfig
	)

	if sc.Perturb, err = i.makePerturber(LowerToPerturber(i.Observers)); err != nil {
		return nil, fmt.Errorf("when making planner: %w", err)
	}
	if sc.Fuzz, err = i.makeFuzzer(bobs); err != nil {
		return nil, fmt.Errorf("when making fuzzer config: %w", err)
	}
	if sc.Lift, err = i.makeLifter(bobs); err != nil {
		return nil, fmt.Errorf("when making lifter config: %w", err)
	}
	if sc.Invoke, err = i.makeInvoker(cobs, LowerToMach(i.Observers)); err != nil {
		return nil, fmt.Errorf("when making machine invoker: %w", err)
	}
	if sc.Analyser, err = i.makeAnalyser(LowerToAnalyser(i.Observers), LowerToSaver(i.Observers)); err != nil {
		return nil, fmt.Errorf("when making analysis: %w", err)
	}
	return &sc, nil
}

func (i *Instance) makeAnalyser(aobs []analyser.Observer, sobs []saver.Observer) (*analyser.Analyser, error) {
	return analyser.New(
		analyser.ObserveWith(aobs...),
		analyser.ObserveSaveWith(sobs...),
		analyser.Analysis(
			analysis.WithWorkerCount(10), // TODO(@MattWindsor91): get this from somewhere
			analysis.WithFilters(i.Filters),
		),
		analyser.SaveToPathset(&i.Machine.Pathset.Saved),
	)
}

func (i *Instance) makePerturber(obs []perturber.Observer) (*perturber.Perturber, error) {
	return perturber.New(
		i.Env.CInspector,
		perturber.ObserveWith(obs...),
		perturber.OverrideQuantities(i.Machine.Quantities.Perturb),
		perturber.UseFullCompilerIDs(true),
	)
}

func (i *Instance) makeFuzzer(obs []builder.Observer) (*fuzzer.Fuzzer, error) {
	return fuzzer.New(
		i.Env.Fuzzer,
		fuzzer.NewPathset(i.Machine.Pathset.Scratch.DirFuzz),
		fuzzer.ObserveWith(obs...),
		fuzzer.OverrideQuantities(i.Machine.Quantities.Fuzz),
		fuzzer.UseConfig(i.FuzzerConfig),
	)
}

func (i *Instance) makeLifter(obs []builder.Observer) (*lifter.Lifter, error) {
	return lifter.New(
		i.Env.BResolver,
		lifter.NewPathset(i.Machine.Pathset.Scratch.DirLift),
		lifter.ObserveWith(obs...),
	)
}

func (i *Instance) makeInvoker(cobs []copier.Observer, mobs []observer2.Observer) (*invoker.Invoker, error) {
	// Unlike the single-shot, we don't late-bind the factory using the plan.  This is because we've already
	// got the machine configuration without it.
	f, err := runner.FactoryFromRemoteConfig(i.SSHConfig, i.Machine.Config.SSH)
	if err != nil {
		return nil, err
	}
	return invoker.New(i.Machine.Pathset.Scratch.DirRun,
		f,
		invoker.ObserveCopiesWith(cobs...),
		invoker.ObserveMachWith(mobs...),
		// As above, there is no loading of quantities using the plan, as we already know which machine the plan is
		// targeting without consulting the plan.
		invoker.OverrideBaseQuantities(i.Machine.Quantities.Mach),
	)
}

func (i *Instance) cleanUpCycle() error {
	return iohelp.Rmdirs(i.Machine.Pathset.Scratch.Dirs()...)
}
