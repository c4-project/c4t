// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package singleobs

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/model/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Logger lifts a Logger to an observer.
type Logger log.Logger

// OnBuildStart does nothing.
func (l *Logger) OnBuildStart(builder.Manifest) {}

// OnBuildRequest logs failed compile and run results.
func (l *Logger) OnBuildRequest(r builder.Request) {
	switch {
	case r.Compile != nil && r.Compile.Result.Status != subject.StatusOk:
		(*log.Logger)(l).Printf("subject %q on compiler %q: %s", r.Name, r.Compile.CompilerID.String(), r.Compile.Result.Status)
	case r.Run != nil && r.Run.Result.Status != subject.StatusOk:
		(*log.Logger)(l).Printf("subject %q on compiler %q: %s", r.Name, r.Run.CompilerID.String(), r.Run.Result.Status)
	}
}

// OnBuildFinish does nothing.
func (l *Logger) OnBuildFinish() {}

// OnCompilerPlanStart briefly logs a compiler start.
func (l *Logger) OnCompilerPlanStart(ncompilers int) {
	(*log.Logger)(l).Printf("planning %d compiler(s)...\n", ncompilers)
}

// OnCompilerPlan does nothing.
func (l *Logger) OnCompilerPlan(_ compiler.Named) {}

// OnCompilerPlanFinish does nothing.
func (l *Logger) OnCompilerPlanFinish() {}
