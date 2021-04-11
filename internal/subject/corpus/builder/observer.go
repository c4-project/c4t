// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package builder

import (
	"github.com/c4-project/c4t/internal/observing"
)

// Observer is the interface for things that observe a builder.
type Observer interface {
	// OnBuild sends a builder observation message.
	OnBuild(Message)
}

//go:generate mockery --name=Observer

// Message is the type of builder observation messages.
type Message struct {
	observing.Batch

	// Manifest carries the name of the subject being (re-)built, if we're on a build-start.
	Name string `json:"name,omitempty"`

	// Request carries a builder request, if we're on a build-step.
	Request *Request `json:"request,omitempty"`
}

// OnBuild sends an OnBuild message to each observer in obs.
func OnBuild(m Message, obs ...Observer) {
	for _, o := range obs {
		o.OnBuild(m)
	}
}

// StartMessage creates an build-start message using manifest m.
func StartMessage(m Manifest) Message {
	return Message{Batch: observing.NewBatchStart(m.NReqs), Name: m.Name}
}

// StepMessage creates a build-step message for step i and request r.
func StepMessage(i int, r Request) Message {
	return Message{Batch: observing.NewBatchStep(i), Request: &r}
}

// EndMessage creates a build-end message.
func EndMessage() Message {
	return Message{Batch: observing.NewBatchEnd()}
}
