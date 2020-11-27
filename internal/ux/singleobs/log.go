// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package singleobs

import (
	"log"

	"github.com/MattWindsor91/c4t/internal/director"

	"github.com/MattWindsor91/c4t/internal/coverage"

	"github.com/MattWindsor91/c4t/internal/plan/analysis"
	"github.com/MattWindsor91/c4t/internal/stage/analyser/saver"

	"github.com/MattWindsor91/c4t/internal/stage/mach/observer"

	"github.com/MattWindsor91/c4t/internal/stage/planner"

	"github.com/MattWindsor91/c4t/internal/stage/perturber"

	"github.com/MattWindsor91/c4t/internal/copier"

	"github.com/MattWindsor91/c4t/internal/observing"

	"github.com/MattWindsor91/c4t/internal/subject/status"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/c4t/internal/model/service/compiler"

	"github.com/MattWindsor91/c4t/internal/subject/corpus/builder"
)

// Logger lifts a Logger to an observer.
type Logger log.Logger

// OnBuild logs build messages.
func (l *Logger) OnBuild(b builder.Message) {
	switch b.Kind {
	case observing.BatchStep:
		l.onBuildRequest(b.Request)
	}
}

// OnBuildRequest logs failed compile and run results.
func (l *Logger) onBuildRequest(r *builder.Request) {
	switch {
	case r.Compile != nil && r.Compile.Result.Status != status.Ok:
		(*log.Logger)(l).Printf("subject %q on compiler %q: %s", r.Name, r.Compile.CompilerID.String(), r.Compile.Result.Status)
	case r.Run != nil && r.Run.Result.Status != status.Ok:
		(*log.Logger)(l).Printf("subject %q on compiler %q: %s", r.Name, r.Run.CompilerID.String(), r.Run.Result.Status)
	}
}

// OnCompilerConfig logs compiler config messages.
func (l *Logger) OnCompilerConfig(m compiler.Message) {
	switch m.Kind {
	case observing.BatchStart:
		l.onCompilerPlanStart(m.Num)
	case observing.BatchStep:
		l.onCompilerPlan(*m.Configuration)
	}
}

// onCompilerPlanStart briefly logs a compiler start.
func (l *Logger) onCompilerPlanStart(ncompilers int) {
	(*log.Logger)(l).Printf("planning %d compiler(s)...\n", ncompilers)
}

// onCompilerPlan logs a compiler plan.
func (l *Logger) onCompilerPlan(nc compiler.Named) {
	(*log.Logger)(l).Printf("compiler %s:\n", nc.ID)
	if nc.SelectedOpt != nil {
		(*log.Logger)(l).Printf(" - opt: %q:\n", nc.SelectedOpt.Name)
	}
	if !ystring.IsBlank(nc.SelectedMOpt) {
		(*log.Logger)(l).Printf(" - m/opt: %q:\n", nc.SelectedMOpt)
	}
}

// OnCopy logs build messages.
func (l *Logger) OnCopy(c copier.Message) {
	switch c.Kind {
	case observing.BatchStart:
		l.onCopyStart(c.Num)
	}
}

// onCopyStart briefly logs a file-copy start.
func (l *Logger) onCopyStart(nfiles int) {
	(*log.Logger)(l).Printf("copying %d files...\n", nfiles)
}

// OnPerturb logs perturb messages.
func (l *Logger) OnPerturb(m perturber.Message) {
	switch m.Kind {
	case perturber.KindStart:
		(*log.Logger)(l).Printf("perturbing plan...\n")
		m.Quantities.Log((*log.Logger)(l))
	case perturber.KindSeedChanged:
		(*log.Logger)(l).Printf("- seed is now %d\n", m.Seed)
	case perturber.KindRandomisingOpts:
		(*log.Logger)(l).Printf("- randomising compiler options...\n")
	case perturber.KindSamplingCorpus:
		(*log.Logger)(l).Printf("- sampling corpus...\n")
	}
}

// OnPlan logs plan messages.
func (l *Logger) OnPlan(m planner.Message) {
	switch m.Kind {
	case planner.KindStart:
		(*log.Logger)(l).Printf("planning...\n")
		m.Quantities.Log((*log.Logger)(l))
	case planner.KindPlanningBackend:
		(*log.Logger)(l).Printf("- probing backend...\n")
	case planner.KindPlanningCompilers:
		(*log.Logger)(l).Printf("- probing compilers on machine %s...\n", m.MachineID)
	case planner.KindPlanningCorpus:
		(*log.Logger)(l).Printf("- probing corpus...\n")
	}
}

// OnMachineNodeAction logs node messages.
func (l *Logger) OnMachineNodeAction(m observer.Message) {
	switch m.Kind {
	case observer.KindCompileStart:
		(*log.Logger)(l).Printf("compiling...\n")
		m.Quantities.Compiler.Log((*log.Logger)(l))
	case observer.KindRunStart:
		(*log.Logger)(l).Printf("running...\n")
		m.Quantities.Runner.Log((*log.Logger)(l))
	}
}

// OnCycle does nothing, for now.
func (l *Logger) OnCycle(director.CycleMessage) {}

// OnInstanceClose does nothing.
func (l *Logger) OnInstanceClose() {}

// OnAnalysis does nothing, for now.
func (l *Logger) OnAnalysis(analysis.Analysis) {}

// OnArchive does nothing, for now.
func (l *Logger) OnArchive(saver.ArchiveMessage) {}

// OnCoverageRun logs information about a coverage run in progress according to rm.
func (l *Logger) OnCoverageRun(rm coverage.RunMessage) {
	switch rm.Kind {
	case observing.BatchStart:
		l.onCoverageRunStart(rm.ProfileName, rm.Num)
	case observing.BatchStep:
		l.onCoverageRunStep(rm.ProfileName, rm.Num)
	case observing.BatchEnd:
		l.onCoverageRunEnd(rm.ProfileName)
	}
}

func (l *Logger) onCoverageRunStart(name string, num int) {
	(*log.Logger)(l).Printf("starting coverage profile %s (%d runs)...\n", name, num)
}

func (l *Logger) onCoverageRunStep(name string, num int) {
	(*log.Logger)(l).Printf("- coverage profile %s: run %d\n", name, num)
}

func (l *Logger) onCoverageRunEnd(name string) {
	(*log.Logger)(l).Printf("finished coverage profile %s\n", name)
}
