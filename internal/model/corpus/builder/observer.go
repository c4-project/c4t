// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package builder

import (
	"github.com/MattWindsor91/act-tester/internal/observing"
)

// Observer is the interface for things that observe a builder.
type Observer interface {
	// OnBuild sends a builder observation message.
	OnBuild(Message)
}

//go:generate mockery -name Observer

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

// OnBuildStart sends an OnBuildStart message to each observer in obs.
func OnBuildStart(m Manifest, obs ...Observer) {
	OnBuild(Message{Batch: observing.NewBatchStart(m.NReqs), Name: m.Name}, obs...)
}

// OnBuildRequest sends an OnBuildRequest message to each observer in obs.
func OnBuildRequest(i int, r Request, obs ...Observer) {
	OnBuild(Message{Batch: observing.NewBatchStep(i), Request: &r}, obs...)
}

// OnBuildStart sends an OnBuildFinish message to each observer in obs.
func OnBuildFinish(obs ...Observer) {
	OnBuild(Message{Batch: observing.NewBatchEnd()}, obs...)
}
