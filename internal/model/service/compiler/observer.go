// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"github.com/c4-project/c4t/internal/observing"
)

// This interface is modelled after the corpus equivalent, builder.Observer, and should probably follow any major
// changes in it.

// Observer is the interface for things that observe the production of compiler configurations.
type Observer interface {
	// OnCompilerConfig sends a compiler configuration observation message.
	OnCompilerConfig(Message)
}

//go:generate mockery --name=Observer

// Message is the type of builder observation messages.
type Message struct {
	observing.Batch
	// Configuration carries a named compiler configuration, if we're on a build-step.
	Configuration *Named `json:"configuration,omitempty"`
}

// OnCompilerConfig sends an OnCompilerConfig message to each observer in obs.
func OnCompilerConfig(m Message, obs ...Observer) {
	for _, o := range obs {
		o.OnCompilerConfig(m)
	}
}

// OnCompilerConfigStart sends a compiler config-start message to each observer in obs.
func OnCompilerConfigStart(nCompilers int, obs ...Observer) {
	OnCompilerConfig(Message{Batch: observing.NewBatchStart(nCompilers)}, obs...)
}

// OnCompilerConfigStep sends an compiler config-step message to each observer in obs.
func OnCompilerConfigStep(i int, c Named, obs ...Observer) {
	OnCompilerConfig(Message{Batch: observing.NewBatchStep(i), Configuration: &c}, obs...)
}

// OnCompilerConfigEnd sends an compiler config-end message to each observer in obs.
func OnCompilerConfigEnd(obs ...Observer) {
	OnCompilerConfig(Message{Batch: observing.NewBatchEnd()}, obs...)
}
