// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package observer defines interfaces and basic implementations of the director's 'observer' pattern.
package director

import (
	"github.com/MattWindsor91/c4t/internal/director/pathset"
	"github.com/MattWindsor91/c4t/internal/quantity"

	"github.com/MattWindsor91/c4t/internal/copier"

	mach "github.com/MattWindsor91/c4t/internal/stage/mach/observer"

	"github.com/MattWindsor91/c4t/internal/stage/perturber"

	"github.com/MattWindsor91/c4t/internal/stage/analyser/saver"

	"github.com/MattWindsor91/c4t/internal/stage/analyser"

	"github.com/MattWindsor91/c4t/internal/machine"

	"github.com/MattWindsor91/c4t/internal/stage/planner"

	"github.com/MattWindsor91/c4t/internal/model/id"
	"github.com/MattWindsor91/c4t/internal/subject/corpus/builder"
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
	// OnCycle observes a message relating to this instance's current cycle.
	OnCycle(m CycleMessage)

	// OnInstanceClose observes that the instance this observer is observing has closed.
	// This gives the observer the opportunity to free any resources.
	OnInstanceClose()

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

// OnInstanceClose sends OnInstanceClose to each observer in obs.
func OnInstanceClose(obs ...InstanceObserver) {
	for _, o := range obs {
		o.OnInstanceClose()
	}
}

type CycleMessage struct {
	// Cycle names the cycle on which this message is occurring.
	Cycle Cycle
	// Kind gives the kind of message.
	Kind CycleMessageKind
	// Err holds the error if Kind is CycleError.
	Err error
}

// CycleMessageKind is the enumeration of kinds of cycle message.
type CycleMessageKind uint8

const (
	// CycleStart denotes the start of a cycle.
	// Future messages from an InstanceObserver should be ascribed to this cycle, until another CycleStart.
	CycleStart CycleMessageKind = iota
	// CycleError denotes a message carrying an error from a cycle.
	// Errors in cycles generally cause the cycle to restart, maybe with backoff.
	CycleError
)

// CycleStartMessage constructs a CycleStart message with cycle c.
func CycleStartMessage(c Cycle) CycleMessage {
	return CycleMessage{Cycle: c, Kind: CycleStart}
}

// CycleErrorMessage constructs a CycleError message with cycle c and error err.
func CycleErrorMessage(c Cycle, err error) CycleMessage {
	return CycleMessage{Cycle: c, Kind: CycleError, Err: err}
}

// OnPrepare sends OnPrepare to every observer in obs.
func OnPrepare(qs quantity.RootSet, ps pathset.Pathset, obs ...Observer) {
	for _, o := range obs {
		o.OnPrepare(qs, ps)
	}
}

// OnCycle sends a cycle message to every instance observer in obs.
func OnCycle(m CycleMessage, obs ...InstanceObserver) {
	for _, o := range obs {
		o.OnCycle(m)
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
