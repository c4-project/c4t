// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/cheggaaa/pb/v3"
)

// SilentObserver is an observer that does nothing.
type SilentObserver struct{}

// OnStart does nothing.
func (s SilentObserver) OnStart(int) {
}

// OnAdd does nothing.
func (s SilentObserver) OnAdd(string) {
}

// OnCompile does nothing.
func (s SilentObserver) OnCompile(string, model.ID) {
}

// OnHarness does nothing.
func (s SilentObserver) OnHarness(string, model.ID) {
}

// OnRun does nothing.
func (s SilentObserver) OnRun(string, model.ID) {
}

// OnFinish does nothing.
func (s SilentObserver) OnFinish() {
}

// BuilderObserver is the interface for things that observe a builder.
type BuilderObserver interface {
	// OnStart executes when the builder starts processing.
	OnStart(nreqs int)

	// OnAdd executes when the subject sname is added to the corpus.
	OnAdd(sname string)

	// OnCompile executes when the subject sname is compiled with compiler cid.
	OnCompile(sname string, cid model.ID)

	// OnHarness executes when the subject sname is lifted to a harness for architecture arch.
	OnHarness(sname string, arch model.ID)

	// OnRun executes when the subject sname's compilation over cid is run.
	OnRun(sname string, cid model.ID)

	// OnFinish executes when the builder stops processing.
	OnFinish()
}

// PbObserver is a builder observer that uses a progress bar.
type PbObserver struct {
	bar *pb.ProgressBar
}

func (p *PbObserver) OnStart(nreqs int) {
	p.bar = pb.StartNew(nreqs)
}

func (p *PbObserver) OnAdd(string) {
	p.bar.Increment()
}

func (p *PbObserver) OnCompile(string, model.ID) {
	p.bar.Increment()
}

func (p *PbObserver) OnHarness(string, model.ID) {
	p.bar.Increment()
}

func (p *PbObserver) OnRun(string, model.ID) {
	p.bar.Increment()
}

func (p *PbObserver) OnFinish() {
	p.bar.Finish()
}
