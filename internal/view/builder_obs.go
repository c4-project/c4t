// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package view

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/cheggaaa/pb/v3"
)

// PbObserver is a builder observer that uses a progress bar.
type PbObserver struct {
	bar *pb.ProgressBar
}

// NewPbObserver creates a new observer using logger l to announce unusual cases.
func NewPbObserver() *PbObserver {
	return &PbObserver{}
}

func (p *PbObserver) OnBuildStart(m builder.Manifest) {
	p.bar = pb.StartNew(m.NReqs)
}

func (p *PbObserver) OnBuildRequest(builder.Request) {
	if p.bar != nil {
		p.bar.Increment()
	}
}

func (p *PbObserver) OnBuildFinish() {
	if p.bar != nil {
		p.bar.Finish()
	}
}

// LogObserver lifts a Logger to an observer.
type LogObserver log.Logger

// OnBuildStart does nothing.
func (l *LogObserver) OnBuildStart(builder.Manifest) {}

// OnBuildRequest logs failed compile and run results.
func (l *LogObserver) OnBuildRequest(r builder.Request) {
	switch {
	case r.Compile != nil && !r.Compile.Result.Success:
		(*log.Logger)(l).Printf("subject %q on compiler %q: compilation failed", r.Name, r.Compile.CompilerID.String())
	case r.Run != nil && r.Run.Result.Status != subject.StatusOk:
		(*log.Logger)(l).Printf("subject %q on compiler %q: %s", r.Name, r.Run.CompilerID.String(), r.Run.Result.Status.String())
	}
}

// OnBuildFinish does nothing.
func (l *LogObserver) OnBuildFinish() {}

// Observers builds a list of observers suitable for single-shot act-tester binaries.
func Observers(l *log.Logger) []builder.Observer {
	return []builder.Observer{
		NewPbObserver(),
		(*LogObserver)(l),
	}
}
