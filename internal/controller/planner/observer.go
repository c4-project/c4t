// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import "github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

// Observer groups all of the disparate observer interfaces that make up an ObserverSet.
// Its main purpose is to let all of those interfaces be implemented by one slice.
type Observer interface {
	builder.Observer
}

// CompilerObserver
type CompilerObserver interface {
}

// ObserverSet groups the various observers used by a planner.
type ObserverSet struct {
	// Corpus contains the corpus-builder observers to be used when building out a plan.
	Corpus []builder.Observer
}

// NewObserverSet creates an ObserverSet using the given observers obs in all roles.
func NewObserverSet(obs ...Observer) ObserverSet {
	lobs := len(obs)
	oset := ObserverSet{
		Corpus: make([]builder.Observer, lobs),
	}
	for i, o := range obs {
		oset.Corpus[i] = o
	}
	return oset
}
