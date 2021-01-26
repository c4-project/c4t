// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"context"

	"github.com/c4-project/c4t/internal/copier"
	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/plan/analysis"
	"github.com/c4-project/c4t/internal/stage/analyser/saver"
	"github.com/c4-project/c4t/internal/stage/mach/observer"
	"github.com/c4-project/c4t/internal/stage/perturber"
	"github.com/c4-project/c4t/internal/subject/corpus/builder"

	"github.com/c4-project/c4t/internal/mutation"
)

// TODO(@MattWindsor91): ideally, this would be decoupled from the director.

func (i *Instance) prepareMutation(ctx context.Context) error {
	if !i.Machine.InitialPlan.IsMutationTest() {
		return nil
	}
	a, err := mutation.NewAutomator(i.Machine.InitialPlan.Mutation.Auto)
	if err != nil {
		return err
	}
	i.mutantCh = a.MutantCh()
	i.Observers = append(i.Observers, killObserver(a.KillCh()))
	go a.Run(ctx)
	return nil
}

func (i *Instance) handleMutantChange(m mutation.Mutant) {
	i.Machine.InitialPlan.SetMutant(m)
	OnInstance(InstanceMutantMessage(m), i.Observers...)
}

// killObserver is the kill channel of a mutation automator adapted to observe instances.
type killObserver chan<- struct{}

// OnCycle does nothing.
func (k killObserver) OnCycle(CycleMessage) {
}

// OnInstance closes the kill channel if the instance is closing.
func (k killObserver) OnInstance(m InstanceMessage) {
	if m.Kind == KindInstanceClosed {
		close(k)
	}
}

// OnAnalysis fires a kill signal if the analysis a suggests the current mutant was killed.
func (k killObserver) OnAnalysis(a analysis.Analysis) {
	if a.Mutation.HasKills() {
		k <- struct{}{}
	}
}

// OnArchive does nothing.
func (k killObserver) OnArchive(saver.ArchiveMessage) {}

// OnCompilerConfig does nothing.
func (k killObserver) OnCompilerConfig(compiler.Message) {}

// OnBuild does nothing.
func (k killObserver) OnBuild(builder.Message) {}

// OnPerturb does nothing.
func (k killObserver) OnPerturb(perturber.Message) {}

// OnCopy does nothing.
func (k killObserver) OnCopy(copier.Message) {}

// OnMachineNodeAction does nothing.
func (k killObserver) OnMachineNodeAction(observer.Message) {}
