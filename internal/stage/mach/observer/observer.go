// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package observer contains interfaces and types for machine node observers.
package observer

import (
	"github.com/c4-project/c4t/internal/quantity"
	"github.com/c4-project/c4t/internal/subject/corpus/builder"
)

// Observer is the interface of anything that observes machine node behaviour.
type Observer interface {
	// This allows observing the corpus builds for both compilation and running.
	builder.Observer

	// OnMachineNodeAction observes a machine node action, usually forwarded from a remote node.
	OnMachineNodeAction(Message)
}

//go:generate mockery --name=Observer

// Message is the type of messages from a machine node observer.
type Message struct {
	// Kind is the kind of message.
	Kind MessageKind
	// Quantities contains, depending on Kind, various pieces of quantity information about the machine node.
	Quantities quantity.MachNodeSet
}

// MessageKind is the enumeration of possible messages from a machine node observer.
type MessageKind uint8

const (
	// The machine node is starting its compile phase.
	// Quantities.Compile contains the compiler quantities.
	KindCompileStart MessageKind = iota
	// The machine node is starting its run phase.
	// Quantities.Run contains the runner quantities.
	KindRunStart
)

// OnMachineNodeAction distributes m to every observer in obs.
func OnMachineNodeAction(m Message, obs ...Observer) {
	for _, o := range obs {
		o.OnMachineNodeAction(m)
	}
}

// OnCompileStart sends a compile start message to every observer in obs, containing the quantity set qs.
func OnCompileStart(qs quantity.BatchSet, obs ...Observer) {
	OnMachineNodeAction(Message{Kind: KindCompileStart, Quantities: quantity.MachNodeSet{Compiler: qs}}, obs...)
}

// OnRunStart sends a run start message to every observer in obs, containing the quantity set qs.
func OnRunStart(qs quantity.BatchSet, obs ...Observer) {
	OnMachineNodeAction(Message{Kind: KindRunStart, Quantities: quantity.MachNodeSet{Runner: qs}}, obs...)
}

// LowerToBuilder lowers each observer in obs to a corpus builder observer.
func LowerToBuilder(obs ...Observer) []builder.Observer {
	bobs := make([]builder.Observer, len(obs))
	for i, o := range obs {
		bobs[i] = o
	}
	return bobs
}
