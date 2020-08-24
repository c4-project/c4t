// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package singleobs

import (
	copy2 "github.com/MattWindsor91/act-tester/internal/copier"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/MattWindsor91/act-tester/internal/observing"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer"
	"github.com/MattWindsor91/act-tester/internal/stage/perturber"
	"github.com/MattWindsor91/act-tester/internal/stage/planner"

	"github.com/cheggaaa/pb/v3"
)

// Bar is a builder observer that uses a progress bar.
type Bar struct {
	bar *pb.ProgressBar
}

// NewBar creates a new observer.
func NewBar() *Bar {
	return &Bar{}
}

// OnBuild observes a corpus build using a progress bar.
func (p *Bar) OnBuild(m builder.Message) {
	p.onBatch(m.Batch)
}

// OnCompilerConfig observes a compiler configuration using a progress bar.
func (p *Bar) OnCompilerConfig(m compiler.Message) {
	p.onBatch(m.Batch)
}

// OnCopy observes a file copy using a progress bar.
func (p *Bar) OnCopy(m copy2.Message) {
	p.onBatch(m.Batch)
}

// OnPerturb does nothing.
func (p *Bar) OnPerturb(perturber.Message) {}

// OnPlan does nothing.
func (p *Bar) OnPlan(planner.Message) {}

// OnMachineNodeAction does nothing.
func (p *Bar) OnMachineNodeAction(observer.Message) {}

func (p *Bar) onBatch(m observing.Batch) {
	switch m.Kind {
	case observing.BatchStart:
		p.start(m.Num)
	case observing.BatchStep:
		p.step()
	case observing.BatchEnd:
		p.finish()
	}
}

func (p *Bar) start(n int) {
	p.bar = pb.StartNew(n)
}

func (p *Bar) step() {
	if p.bar != nil {
		p.bar.Increment()
	}
}

func (p *Bar) finish() {
	if p.bar != nil {
		p.bar.Finish()
	}
}
