// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package observer defines interfaces and basic implementations of the director's 'observer' pattern.
package director

import (
	"github.com/c4-project/c4t/internal/director/pathset"
	"github.com/c4-project/c4t/internal/mutation"
	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/copier"

	mach "github.com/c4-project/c4t/internal/stage/mach/observer"

	"github.com/c4-project/c4t/internal/stage/perturber"

	"github.com/c4-project/c4t/internal/stage/analyser/saver"

	"github.com/c4-project/c4t/internal/stage/analyser"

	"github.com/c4-project/c4t/internal/machine"

	"github.com/c4-project/c4t/internal/stage/planner"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/subject/corpus/builder"
)

// Observer is an interface for types that implement multi-machine test progress observation.
type Observer interface {
	// Observers can observe machine configuration.
	machine.Observer

	// Observers can observe planner operations.
	planner.Observer

	// Observers can observe director preparations.
	PrepareObserver

	// Instance gets a sub-observer for the machine with ID id.
	// It can fail if no such observer is available.
	Instance(id id.ID) (InstanceObserver, error)
}

// PrepareObserver is an interface for types that observer director preparations.
type PrepareObserver interface {
	// TODO(@MattWindsor91): this can go away if we flatten things into single, auto-forwarded Observers

	// OnPrepare lets the observer know that the director is performing some sort of preparation.
	OnPrepare(message PrepareMessage)
}

// PrepareKind is the enumeration of kinds of PrepareMessage.
type PrepareKind uint8

const (
	// PrepareInstances states that the director is preparing its instance set; NumInstances is set.
	PrepareInstances PrepareKind = iota
	// PrepareQuantities states the director's quantity set; Quantities is set.
	PrepareQuantities
	// PreparePaths states that the director is about to make its top-level paths; Paths is set.
	PreparePaths
)

// PrepareMessage is a message from the director stating some aspect of its pre-experiment preparation.
type PrepareMessage struct {
	// Kind states the kind of message.
	Kind PrepareKind `json:"kind"`

	// NumInstances states, if Kind is PrepareInstances, the number of instances that the director is going to run.
	// This can be used by observers to pre-allocate instance sub-observers.
	NumInstances int `json:"num_instances,omitempty"`
	// Quantities states, if Kind is PrepareQuantities, the quantities the director is going to make.
	Quantities quantity.RootSet `json:"quantities,omitempty"`
	// Paths states, if Kind is PreparePaths, where the director is going to make its top-level paths.
	Paths pathset.Pathset `json:"paths,omitempty"`
}

// PrepareInstancesMessage creates a PrepareMessage with kind PrepareInstances and instance count ninst.
func PrepareInstancesMessage(ninst int) PrepareMessage {
	return PrepareMessage{Kind: PrepareInstances, NumInstances: ninst}
}

// PrepareQuantitiesMessage creates a PrepareMessage with kind PrepareQuantities and quantity set qs.
func PrepareQuantitiesMessage(qs quantity.RootSet) PrepareMessage {
	return PrepareMessage{Kind: PrepareQuantities, Quantities: qs}
}

// PreparePathsMessage creates a PrepareMessage with kind PreparePaths and path set ps.
func PreparePathsMessage(ps pathset.Pathset) PrepareMessage {
	return PrepareMessage{Kind: PreparePaths, Paths: ps}
}

// OnPrepare sends OnPrepare to every observer in obs.
func OnPrepare(m PrepareMessage, obs ...PrepareObserver) {
	for _, o := range obs {
		o.OnPrepare(m)
	}
}

// CycleObserver is an interface for types that observe cycles.
//
// This is a separate sub-interface of InstanceObserver because some things implement one, but not the other.
type CycleObserver interface {
	// OnCycle observes a message relating to this instance's current cycle.
	OnCycle(m CycleMessage)
}

// InstanceObserver is an interface for types that observe a director instance.
type InstanceObserver interface {
	// InstanceObserver observers can observe cycles.
	CycleObserver

	// OnInstance observes something about that the instance this observer is observing.
	OnInstance(m InstanceMessage)

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

// InstanceMessage is the type of observer messages pertaining to the control flow of an instance.
type InstanceMessage struct {
	Kind   InstanceMessageKind
	Mutant mutation.Mutant
}

// InstanceMessageKind is the enumeration of kinds of instance message.
type InstanceMessageKind uint8

const (
	// The instance has closed.
	// Observers should free any resources specific to this instance.
	KindInstanceClosed InstanceMessageKind = iota
	// The instance has changed to a new mutant (in Mutant).
	KindInstanceMutant
)

// InstanceClosedMessage constructs an InstanceMessage stating that the instance has closed.
func InstanceClosedMessage() InstanceMessage {
	return InstanceMessage{Kind: KindInstanceClosed}
}

// InstanceMutantMessage constructs an InstanceMessage stating that the instance has changed mutant to m.
func InstanceMutantMessage(m mutation.Mutant) InstanceMessage {
	return InstanceMessage{Kind: KindInstanceMutant, Mutant: m}
}

// OnInstance sends OnInstance to each observer in obs.
func OnInstance(m InstanceMessage, obs ...InstanceObserver) {
	for _, o := range obs {
		o.OnInstance(m)
	}
}

// CycleMessage is the type of observer messages pertaining to the control flow of a specific cycle.
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
	// CycleFinish denotes the successful completion of a cycle.
	CycleFinish
	// CycleError denotes a message carrying an error from a cycle (replacing CycleFinish).
	// Errors in cycles generally cause the cycle to restart, maybe with backoff.
	CycleError
)

// CycleStartMessage constructs a CycleStart message with cycle c.
func CycleStartMessage(c Cycle) CycleMessage {
	return CycleMessage{Cycle: c, Kind: CycleStart}
}

// CycleFinishMessage constructs a CycleFinish message with cycle c.
func CycleFinishMessage(c Cycle) CycleMessage {
	return CycleMessage{Cycle: c, Kind: CycleFinish}
}

// CycleErrorMessage constructs a CycleError message with cycle c and error err.
func CycleErrorMessage(c Cycle, err error) CycleMessage {
	return CycleMessage{Cycle: c, Kind: CycleError, Err: err}
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

// LowerToPrepare lowers a slice of director observers to a slice of prepare observers.
func LowerToPrepare(obs []Observer) []PrepareObserver {
	dobs := make([]PrepareObserver, len(obs))
	for i, o := range obs {
		dobs[i] = o
	}
	return dobs
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
