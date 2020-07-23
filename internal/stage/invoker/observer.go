// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package invoker

import (
	"github.com/MattWindsor91/act-tester/internal/copier"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
)

// Observer is the union of the various interfaces of observers used by invoker.
type Observer interface {
	copier.Observer
	builder.Observer
}

// ObserverSet is a set of observers for use by invoker.
type ObserverSet struct {
	// Copy contains observers that listen for file copies.
	Copy []copier.Observer
	// Corpus contains observers that listen for corpus-building activity on the remote machine.
	Corpus []builder.Observer
}

// NewObserverSet constructs an observer set by using obs in all roles.
func NewObserverSet(obs ...Observer) ObserverSet {
	lobs := len(obs)
	oset := ObserverSet{
		Copy:   make([]copier.Observer, lobs),
		Corpus: make([]builder.Observer, lobs),
	}
	for i, o := range obs {
		oset.Corpus[i] = o
		oset.Copy[i] = o
	}
	return oset
}

// Append appends the observers in os to this set.
func (o *ObserverSet) Append(os ObserverSet) {
	o.Copy = append(o.Copy, os.Copy...)
	o.Corpus = append(o.Corpus, os.Corpus...)
}
