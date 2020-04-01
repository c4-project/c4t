// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package singleobs

import (
	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

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
func (p *Bar) OnBuildStart(m builder.Manifest) {
	p.bar = pb.StartNew(m.NReqs)
}

// OnBuildRequest observes a request on a corpus build using a progress bar.
func (p *Bar) OnBuildRequest(builder.Request) {
	if p.bar != nil {
		p.bar.Increment()
	}
}

// OnBuildFinish observes the end of a corpus build using a progress bar.
func (p *Bar) OnBuildFinish() {
	if p.bar != nil {
		p.bar.Finish()
	}
}
