// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package singleobs

import (
	copy2 "github.com/MattWindsor91/c4t/internal/copier"
	"github.com/MattWindsor91/c4t/internal/director"
	"github.com/MattWindsor91/c4t/internal/model/service/compiler"
	"github.com/MattWindsor91/c4t/internal/observing"
	"github.com/MattWindsor91/c4t/internal/plan/analysis"
	"github.com/MattWindsor91/c4t/internal/stage/analyser/saver"
	"github.com/MattWindsor91/c4t/internal/stage/mach/observer"
	"github.com/MattWindsor91/c4t/internal/stage/perturber"
	"github.com/MattWindsor91/c4t/internal/stage/planner"
	"github.com/MattWindsor91/c4t/internal/subject/corpus/builder"

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

// OnCycle does nothing.
func (p *Bar) OnCycle(director.CycleMessage) {}

// OnInstanceClose does nothing.
func (p *Bar) OnInstanceClose() {}

// OnAnalysis does nothing.
func (p *Bar) OnAnalysis(analysis.Analysis) {}

// OnArchive does nothing.
func (p *Bar) OnArchive(saver.ArchiveMessage) {}

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
