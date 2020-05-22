// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package builder

// Observer is the interface for things that observe a builder.
type Observer interface {
	// OnBuild sends a builder observation message.
	OnBuild(Message)
}

//go:generate mockery -name Observer

// MessageKind is the enumeration of kinds of message.
type MessageKind uint8

const (
	// BuildStart signifies the start of a corpus build.
	BuildStart MessageKind = iota
	// BuildRequest signifies a step in a corpus build.
	BuildRequest
	// BuildFinish signifies the end of a corpus build.
	BuildFinish
)

// Message is the type of builder observation messages.
type Message struct {
	// Kind is the kind of this message.
	Kind MessageKind `json:"kind"`
	// Manifest carries a builder manifest, if we're on a build-start.
	Manifest *Manifest `json:"manifest,omitempty"`

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
	OnBuild(Message{Kind: BuildStart, Manifest: &m}, obs...)
}

// OnBuildRequest sends an OnBuildRequest message to each observer in obs.
func OnBuildRequest(r Request, obs ...Observer) {
	OnBuild(Message{Kind: BuildRequest, Request: &r}, obs...)
}

// OnBuildStart sends an OnBuildFinish message to each observer in obs.
func OnBuildFinish(obs ...Observer) {
	OnBuild(Message{Kind: BuildFinish}, obs...)
}
