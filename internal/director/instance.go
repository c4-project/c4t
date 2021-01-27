// Copyright (c) 2020-2021 C4 Project
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

	"github.com/c4-project/c4t/internal/stage/perturber"

	"github.com/c4-project/c4t/internal/helper/errhelp"

	"github.com/c4-project/c4t/internal/stage/analyser"

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
	FuzzerConfig *fuzzer2.Config

	// mutantCh stores a channel that will receive mutations, if any.
	mutantCh <-chan mutation.Mutant

	// timeoutCh stores the current error cooldown channel, if any.
	// This is refreshed whenever an error occurs.
	timeoutCh <-chan time.Time

	// cycleCh stores the current cycle result channel, if any.
	// This is refreshed whenever a new cycle is launched.
	cycleCh <-chan cycleResult
}

// Run runs this instance's testing loop.
func (i *Instance) Run(ctx context.Context) error {
	err := i.runInner(ctx)
	cerr := i.cleanUp()
	OnInstance(InstanceClosedMessage(), i.Observers...)
	return errhelp.FirstError(err, cerr)
}

// runInner runs the instance's testing loop, less some closedown boilerplate.
func (i *Instance) runInner(ctx context.Context) error {
	if err := i.prepare(ctx); err != nil {
		return err
	}
	return i.mainLoop(ctx)
}

func (i *Instance) prepare(ctx context.Context) error {
	var err error
	if err = i.check(); err != nil {
		return err
	}
	if err = i.Machine.Pathset.Scratch.Prepare(); err != nil {
		return err
	}
	// TODO(@MattWindsor91): move this out of the instance, if possible.
	if err := i.prepareMutation(ctx); err != nil {
		return err
	}
	// This must happen after preparing the mutation config, otherwise the kill channel won't be installed.
	if i.Machine.stages, err = i.makeStages(); err != nil {
		return err
	}

	return nil
}

// cleanUp closes things that should be gracefully closed after an instance terminates.
func (i *Instance) cleanUp() error {
	if i.Machine == nil {
		return nil
	}
	return i.Machine.cleanUp()
}

func (m *Machine) cleanUp() error {
	var err error
	for _, r := range m.stages {
		err = r.Close()
	}
	return err
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

// cycleResult is the type of results from cycle goroutines.
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
		case m := <-i.mutantCh:
			i.handleMutantChange(m)
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
	OnCycle(CycleFinishMessage(res.cycle), i.Observers...)
	i.Machine.cycle++
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
		p:      i.plan(),
		stages: i.Machine.stages,
	}
}

func (i *Instance) plan() *plan.Plan {
	// Important to _copy_ the plan
	pcopy := i.Machine.InitialPlan
	return &pcopy
}

// makeStages constructs a slice of stage runners that will be
func (i *Instance) makeStages() ([]plan.Runner, error) {
	var stages []plan.Runner

	for _, f := range []func() (plan.Runner, error){
		i.makePerturber,
		i.makeFuzzer,
		i.makeLifter,
		i.makeInvoker,
		i.makeAnalyser,
	} {
		s, err := f()
		if err != nil {
			return nil, err
		}
		if s != nil {
			stages = append(stages, s)
		}
	}
	return stages, nil
}

func (i *Instance) makeAnalyser() (plan.Runner, error) {
	return analyser.New(
		analyser.ObserveWith(LowerToAnalyser(i.Observers)...),
		analyser.ObserveSaveWith(LowerToSaver(i.Observers)...),
		analyser.Analysis(
			analysis.WithWorkerCount(10), // TODO(@MattWindsor91): get this from somewhere
			analysis.WithFilters(i.Filters),
		),
		analyser.SaveToPathset(&i.Machine.Pathset.Saved),
	)
}

func (i *Instance) makePerturber() (plan.Runner, error) {
	return perturber.New(
		i.Env.CInspector,
		perturber.ObserveWith(LowerToPerturber(i.Observers)...),
		perturber.OverrideQuantities(i.Machine.Quantities.Perturb),
		perturber.UseFullCompilerIDs(true),
	)
}

// makeFuzzer makes a plan runner for the fuzzer stage.
// If the fuzzer is disabled, this returns nil.
func (i *Instance) makeFuzzer() (plan.Runner, error) {
	if i.FuzzerConfig != nil && i.FuzzerConfig.Disabled {
		return nil, nil
	}

	return fuzzer.New(
		i.Env.Fuzzer,
		fuzzer.NewPathset(i.Machine.Pathset.Scratch.DirFuzz),
		fuzzer.ObserveWith(LowerToBuilder(i.Observers)...),
		fuzzer.OverrideQuantities(i.Machine.Quantities.Fuzz),
		fuzzer.UseConfig(i.FuzzerConfig),
	)
}

func (i *Instance) makeLifter() (plan.Runner, error) {
	return lifter.New(
		i.Env.BResolver,
		lifter.NewPathset(i.Machine.Pathset.Scratch.DirLift),
		lifter.ObserveWith(LowerToBuilder(i.Observers)...),
	)
}

func (i *Instance) makeInvoker() (plan.Runner, error) {
	// Unlike the single-shot, we don't late-bind the factory using the plan.  This is because we've already
	// got the machine configuration without it.
	f, err := runner.FactoryFromRemoteConfig(i.SSHConfig, i.Machine.Config.SSH)
	if err != nil {
		return nil, err
	}
	return invoker.New(i.Machine.Pathset.Scratch.DirRun,
		f,
		invoker.ObserveCopiesWith(LowerToCopy(i.Observers)...),
		invoker.ObserveMachWith(LowerToMach(i.Observers)...),
		// As above, there is no loading of quantities using the plan, as we already know which machine the plan is
		// targeting without consulting the plan.
		invoker.OverrideBaseQuantities(i.Machine.Quantities.Mach),
	)
}

func (i *Instance) cleanUpCycle() error {
	return iohelp.Rmdirs(i.Machine.Pathset.Scratch.Dirs()...)
}
