// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package singleobs

import (
	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/cheggaaa/pb/v3"
)

// Bar is a builder observer that uses a progress bar.
type Bar struct {
	bar *pb.ProgressBar
}

// NewBar creates a new observer using logger l to announce unusual cases.
func NewBar() *Bar {
	return &Bar{}
}

// OnBuildStart observes the start of a corpus build using a progress bar.
func (p *Bar) OnBuild(m builder.Message) {
	switch m.Kind {
	case builder.BuildStart:
		p.start(m.Manifest.NReqs)
	case builder.BuildRequest:
		p.step()
	case builder.BuildFinish:
		p.finish()
	}
}

// OnCompilerPlanStart observes the start of a compiler plan using a progress bar.
func (p *Bar) OnCompilerPlanStart(ncompilers int) {
	p.start(ncompilers)
}

// OnCompilerPlan observes a step of a compiler plan using a progress bar.
func (p *Bar) OnCompilerPlan(compiler.Named) {
	p.step()
}

// OnCompilerPlanFinish observes the end of a compiler plan using a progress bar.
func (p *Bar) OnCompilerPlanFinish() {
	p.finish()
}

// OnCopyStart observes the start of a file copy using a progress bar.
func (p *Bar) OnCopyStart(nfiles int) {
	p.start(nfiles)
}

// OnCopy observes a step of a file copy using a progress bar.
func (p *Bar) OnCopy(_, _ string) {
	p.step()
}

// OnCopyFinish observes the end of a file copy using a progress bar.
func (p *Bar) OnCopyFinish() {
	p.finish()
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
