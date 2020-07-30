// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package perturber

import (
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
)

// Observer is the type of observers for the perturber.
type Observer interface {
	compiler.Observer
	builder.Observer

	// OnPerturb is sent when the perturber is doing something new.
	OnPerturb(m Message)
}

// Kind is the enumeration of kinds of perturber message.
type Kind uint8

const (
	// The perturber is starting.  Quantities points to the perturber's quantity set.
	KindStart Kind = iota
	// The perturber is now sampling the corpus.
	KindSampleCorpus
	// The perturber is now randomising the compiler optimisations.
	KindRandomiseOpts
)

// Message is the type of messages sent through OnPerturb.
type Message struct {
	// Kind is the kind of message being sent.
	Kind Kind

	// Quantities points to the quantity set on start messages.
	Quantities *QuantitySet
}

// OnPerturb sends a perturb message m to each observer in obs.
func OnPerturb(m Message, obs ...Observer) {
	for _, o := range obs {
		o.OnPerturb(m)
	}
}
