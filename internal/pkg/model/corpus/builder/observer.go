// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package builder

// Observer is the interface for things that observe a builder.
type Observer interface {
	// OnStart executes when the builder starts processing.
	OnStart(Manifest)

	// OnRequest executes when a corpus builder request is processed.
	OnRequest(Request)

	// OnFinish executes when the builder stops processing.
	OnFinish()
}

// SilentObserver is an observer that does nothing.
type SilentObserver struct{}

// OnStart does nothing.
func (s SilentObserver) OnStart(Manifest) {
}

// OnUpdate does nothing.
func (s SilentObserver) OnRequest(Request) {
}

// OnFinish does nothing.
func (s SilentObserver) OnFinish() {
}
