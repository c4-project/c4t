// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/cheggaaa/pb/v3"
)

// BuilderConfig is a configuration for a Builder.
type BuilderConfig struct {
	// Init is the initial corpus.
	// If nil, the Builder starts with a new corpus with capacity equal to NReqs.
	// Otherwise, it copies this corpus.
	Init Corpus

	// NReqs is the number of expected requests to be made to the Builder.
	// The builder will finish listening for requests when this target is reached.
	NReqs int

	// Obs is the observer to notify as the builder performs various tasks.
	Obs BuilderObserver
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
func (s SilentObserver) OnCompile(string, model.MachQualID) {
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
	OnCompile(sname string, cid model.MachQualID)

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

func (p *PbObserver) OnCompile(string, model.MachQualID) {
	p.bar.Increment()
}

func (p *PbObserver) OnFinish() {
	p.bar.Finish()
}
