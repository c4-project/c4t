// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package rmach

import (
	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/remote"
)

// Observer is the union of the various interfaces of observers used by rmach.
type Observer interface {
	remote.CopyObserver
	builder.Observer
}

// ObserverSet is a set of observers for use by rmach.
type ObserverSet struct {
	// Copy contains observers that listen for file copies.
	Copy []remote.CopyObserver
	// Corpus contains observers that listen for corpus-building activity on the remote machine.
	Corpus []builder.Observer
}

// NewObserverSet constructs an observer set by using obs in all roles.
func NewObserverSet(obs ...Observer) ObserverSet {
	lobs := len(obs)
	oset := ObserverSet{
		Copy:   make([]remote.CopyObserver, lobs),
		Corpus: make([]builder.Observer, lobs),
	}
	for i, o := range obs {
		oset.Corpus[i] = o
		oset.Copy[i] = o
	}
	return oset
}
