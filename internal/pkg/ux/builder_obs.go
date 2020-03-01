// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package ux

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
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

func (p *PbObserver) OnStart(nreqs int) {
	p.bar = pb.StartNew(nreqs)
}

func (p *PbObserver) OnAdd(string) {
	if p.bar != nil {
		p.bar.Increment()
	}
}

func (p *PbObserver) OnCompile(name string, cid model.ID, success bool) {
	if !success {
		p.l.Printf("subject %q on compiler %q: compilation failed", name, cid.String())
	}
	if p.bar != nil {
		p.bar.Increment()
	}
}

func (p *PbObserver) OnHarness(string, model.ID) {
	if p.bar != nil {
		p.bar.Increment()
	}
}

func (p *PbObserver) OnRun(name string, cid model.ID, s subject.Status) {
	if s != subject.StatusOk {
		p.l.Printf("subject %q on compiler %q: %s", name, cid.String(), s.String())
	}
	if p.bar != nil {
		p.bar.Increment()
	}
}

func (p *PbObserver) OnFinish() {
	if p.bar != nil {
		p.bar.Finish()
	}
}
