// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// BuilderObserver is the interface for things that observe a builder.
type BuilderObserver interface {
	// OnStart executes when the builder starts processing.
	OnStart(nreqs int)

	// OnAdd executes when the subject sname is added to the corpus.
	OnAdd(sname string)

	// OnCompile executes when the subject sname is compiled with compiler cid.
	OnCompile(sname string, cid model.ID, success bool)

	// OnHarness executes when the subject sname is lifted to a harness for architecture arch.
	OnHarness(sname string, arch model.ID)

	// OnRun executes when the subject sname's compilation over cid is run.
	OnRun(sname string, cid model.ID, status subject.Status)

	// OnFinish executes when the builder stops processing.
	OnFinish()
}

// SilentObserver is an observer that does nothing.
type SilentObserver struct{}

// OnStart does nothing.
func (s SilentObserver) OnStart(int) {
}

// OnAdd does nothing.
func (s SilentObserver) OnAdd(string) {
}

// OnCompile does nothing.
func (s SilentObserver) OnCompile(string, model.ID, bool) {
}

// OnHarness does nothing.
func (s SilentObserver) OnHarness(string, model.ID) {
}

// OnRun does nothing.
func (s SilentObserver) OnRun(string, model.ID, subject.Status) {
}

// OnFinish does nothing.
func (s SilentObserver) OnFinish() {
}
