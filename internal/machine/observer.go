// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package machine

// Observer is the observer interface for machine configuration.
type Observer interface {
	// OnMachines sends an observation about the machines defined on a tester.
	OnMachines(m Message)
}

//go:generate mockery --name=Observer

// MessageKind is the enumeration of machine observer message kinds.
type MessageKind uint8

const (
	// MessageStart signifies that this machine message is starting a machine block.
	// Index stores the number of machines being defined.
	MessageStart MessageKind = iota
	// MessageRecord signifies that this machine message is defining a machine.
	// Index stores the index of the machine being defined.
	MessageRecord
	// MessageFinish signifies the end of a machine block.
	MessageFinish
)

// Message outlines a machine information message.
type Message struct {
	// Kind is the kind of message being sent.
	Kind MessageKind
	// Index is the index, or count, depending on the kind.
	Index int
	// Machine, if present, gives information about the machine being described.
	Machine *Named
}

// OnMachines broadcasts the message m to every observer in obs.
func OnMachines(m Message, obs ...Observer) {
	for _, o := range obs {
		o.OnMachines(m)
	}
}

// OnMachinesStart broadcasts the start of a machines block to every observer in obs.
func OnMachinesStart(nmachines int, obs ...Observer) {
	OnMachines(Message{
		Kind:  MessageStart,
		Index: nmachines,
	}, obs...)
}

// OnMachinesRecord broadcasts a record in a machines block to every observer in obs.
func OnMachinesRecord(index int, machine Named, obs ...Observer) {
	OnMachines(Message{
		Kind:    MessageRecord,
		Index:   index,
		Machine: &machine,
	}, obs...)
}

// OnMachinesFinish broadcasts the end of a machines block to every observer in obs.
func OnMachinesFinish(obs ...Observer) {
	OnMachines(Message{Kind: MessageFinish}, obs...)
}
