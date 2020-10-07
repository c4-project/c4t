// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package perturber

import (
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"
)

// Observer is the type of observers for the perturber.
type Observer interface {
	compiler.Observer
	builder.Observer

	// OnPerturb is sent when the perturber is doing something new.
	OnPerturb(m Message)
}

//go:generate mockery --name=Observer

// Kind is the enumeration of kinds of perturber message.
type Kind uint8

const (
	// The perturber is starting.
	// Quantities points to the perturber's quantity set.
	KindStart Kind = iota
	// The perturber has now changed the seed.
	// Seed points to the new seed.
	KindSeedChanged
	// The perturber is now sampling the corpus.
	// The selected corpus will be announced as a series of OnBuild messages.
	KindSamplingCorpus
	// The perturber is now randomising the compiler optimisations.
	// The selected compilers will be announced as a series of OnCompilerConfig messages.
	KindRandomisingOpts
)

// Message is the type of messages sent through OnPerturb.
type Message struct {
	// Kind is the kind of message being sent.
	Kind Kind

	// Quantities points to the quantity set on start messages.
	Quantities *quantity.PerturbSet

	// Seed points to the seed set on seed-changed messages.
	Seed int64
}

// OnPerturb sends a perturb message m to each observer in obs.
func OnPerturb(m Message, obs ...Observer) {
	for _, o := range obs {
		o.OnPerturb(m)
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
