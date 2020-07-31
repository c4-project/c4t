// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
)

// Observer is the type of observers for the planner.
type Observer interface {
	compiler.Observer
	builder.Observer

	// OnPlan is sent when the planner is doing something new.
	OnPlan(m Message)
}

// Kind is the enumeration of kinds of planner message.
type Kind uint8

const (
	// The planner is starting.
	// Quantities points to the planner's quantity set.
	KindStart Kind = iota
	// The planner is now getting information about the backend to use.
	KindPlanningBackend
	// The planner is now getting information about the corpus.
	// The corpus will be announced as a series of OnBuild messages.
	KindPlanningCorpus
	// The planner is now getting information about the compilers for a given machine.
	// MachineID points to the name of the machine.
	// The selected compilers will be announced as a series of OnCompilerConfig messages.
	KindPlanningCompilers
)

// Message is the type of messages sent through OnPlan.
type Message struct {
	// Kind is the kind of message being sent.
	Kind Kind

	// Quantities points to the quantity set on start messages.
	Quantities *QuantitySet

	// MachineID contains the machine identifier in certain messages.
	MachineID id.ID
}

// OnPlan sends a plan message m to each observer in obs.
func OnPlan(m Message, obs ...Observer) {
	for _, o := range obs {
		o.OnPlan(m)
	}
}

func lowerToBuilder(obs []Observer) []builder.Observer {
	cobs := make([]builder.Observer, len(obs))
	for i, o := range obs {
		cobs[i] = o
	}
	return cobs
}

func lowerToCompiler(obs []Observer) []compiler.Observer {
	cobs := make([]compiler.Observer, len(obs))
	for i, o := range obs {
		cobs[i] = o
	}
	return cobs
}
