// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package observer defines interfaces and basic implementations of the director's 'observer' pattern.
package director

import (
	"github.com/MattWindsor91/act-tester/internal/director/pathset"
	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/copier"

	mach "github.com/MattWindsor91/act-tester/internal/stage/mach/observer"

	"github.com/MattWindsor91/act-tester/internal/stage/perturber"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser"

	"github.com/MattWindsor91/act-tester/internal/machine"

	"github.com/MattWindsor91/act-tester/internal/stage/planner"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"
)

// Observer is an interface for types that implement multi-machine test progress observation.
type Observer interface {
	// Observers can observe machine configuration.
	machine.Observer

	// Observers can observe planner operations.
	planner.Observer

	// OnPrepare lets the observer know that the director is preparing to run on pathset ps with quantities qs.
	OnPrepare(qs quantity.RootSet, ps pathset.Pathset)

	// InstanceObserver gets a sub-observer for the machine with ID id.
	// It can fail if no such observer is available.
	Instance(id id.ID) (InstanceObserver, error)
}

// InstanceObserver is an interface for types that observe a director instance.
type InstanceObserver interface {
	// OnIteration lets the observer know that the instance has started a new cycle.
	OnIteration(c Cycle)

	// InstanceObserver observers can observe plan analyses.
	analyser.Observer

	// InstanceObserver observers can observe plan saves.
	saver.Observer

	// InstanceObserver observers can observe perturber operations.
	perturber.Observer

	// InstanceObserver observers can observe file copies.
	copier.Observer

	// InstanceObserver observers can observe machine node actions.
	mach.Observer
}

// OnPrepare sends OnPrepare to every observer in obs.
func OnPrepare(qs quantity.RootSet, ps pathset.Pathset, obs ...Observer) {
	for _, o := range obs {
		o.OnPrepare(qs, ps)
	}
}

// OnIteration sends OnIteration to every instance observer in obs.
func OnIteration(r Cycle, obs ...InstanceObserver) {
	for _, o := range obs {
		o.OnIteration(r)
	}
}

// LowerToPlanner lowers a slice of director observers to a slice of planner observers.
func LowerToPlanner(obs []Observer) []planner.Observer {
	cos := make([]planner.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToMachine lowers a slice of director observers to a slice of machine observers.
func LowerToMachine(obs []Observer) []machine.Observer {
	cos := make([]machine.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToAnalyser lowers a slice of instance observers to a slice of analyser observers.
func LowerToAnalyser(obs []InstanceObserver) []analyser.Observer {
	cos := make([]analyser.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToPerturber lowers a slice of instance observers to a slice of perturber observers.
func LowerToPerturber(obs []InstanceObserver) []perturber.Observer {
	cos := make([]perturber.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToSave lowers a slice of instance observers to a slice of saver observers.
func LowerToSaver(obs []InstanceObserver) []saver.Observer {
	cos := make([]saver.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToBuilder lowers a slice of instance observers to a slice of builder observers.
func LowerToBuilder(obs []InstanceObserver) []builder.Observer {
	cos := make([]builder.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToCopy lowers a slice of instance observers to a slice of copy observers.
func LowerToCopy(obs []InstanceObserver) []copier.Observer {
	cos := make([]copier.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToMach lowers a slice of director observers to a slice of machine node observers.
func LowerToMach(obs []InstanceObserver) []mach.Observer {
	cos := make([]mach.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}
