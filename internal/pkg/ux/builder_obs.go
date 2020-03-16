// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package ux

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
	"github.com/cheggaaa/pb/v3"
)

// PbObserver is a builder observer that uses a progress bar.
type PbObserver struct {
	bar *pb.ProgressBar
	l   *log.Logger
}

// NewPbObserver creates a new observer using logger l to announce unusual cases.
func NewPbObserver(l *log.Logger) *PbObserver {
	return &PbObserver{l: iohelp.EnsureLog(l)}
}

func (p *PbObserver) OnStart(m builder.Manifest) {
	p.bar = pb.StartNew(m.NReqs)
}

func (p *PbObserver) OnRequest(r builder.Request) {
	if p.bar != nil {
		p.bar.Increment()
	}
	switch {
	case r.Compile != nil && !r.Compile.Result.Success:
		p.l.Printf("subject %q on compiler %q: compilation failed", r.Name, r.Compile.CompilerID.String())
	case r.Run != nil && r.Run.Result.Status != subject.StatusOk:
		p.l.Printf("subject %q on compiler %q: %s", r.Name, r.Run.CompilerID.String(), r.Run.Result.Status.String())
	}
}

func (p *PbObserver) OnFinish() {
	if p.bar != nil {
		p.bar.Finish()
	}
}
