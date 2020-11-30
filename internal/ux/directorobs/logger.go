// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package directorobs

import (
	"io"
	"log"

	"github.com/MattWindsor91/c4t/internal/director/pathset"
	"github.com/MattWindsor91/c4t/internal/quantity"

	"github.com/MattWindsor91/c4t/internal/director"

	"github.com/MattWindsor91/c4t/internal/helper/stringhelp"

	"github.com/MattWindsor91/c4t/internal/stage/planner"

	"github.com/MattWindsor91/c4t/internal/observing"

	"github.com/MattWindsor91/c4t/internal/stage/analyser/saver"

	"github.com/MattWindsor91/c4t/internal/stage/analyser/pretty"

	"github.com/MattWindsor91/c4t/internal/machine"

	"github.com/MattWindsor91/c4t/internal/model/service/compiler"

	"github.com/MattWindsor91/c4t/internal/subject/corpus/builder"
)

// TODO(@MattWindsor91): merge this with the singleobs logger?

// Logger is a ForwardHandler that emits logs to a writer when cycles finish up.
type Logger struct {
	// out is the writer to use for logging analyses.
	out io.WriteCloser
	// l is the intermediate logger that sits atop out.
	l *log.Logger
	// aw is the analyser writer used for outputting sourced analyses.
	aw *pretty.Printer
	// compilers holds state for assembling compiler builds.
	compilers map[string][]compiler.Named
}

// OnCompilerConfig (currently) does nothing.
func (j *Logger) OnCompilerConfig(compiler.Message) {
}

// OnBuild (currently) does nothing.
func (j *Logger) OnBuild(builder.Message) {
}

// OnPlan (currently) does nothing.
func (j *Logger) OnPlan(planner.Message) {
}

// NewLogger constructs a new Logger writing into w, using logger flags lflag when logging things.
// The logger takes ownership of w.
func NewLogger(w io.WriteCloser, lflag int) (*Logger, error) {
	aw, err := pretty.NewPrinter(
		pretty.WriteTo(w),
		pretty.ShowCompilers(true),
		pretty.ShowSubjects(true),
	)
	if err != nil {
		return nil, err
	}
	l := &Logger{
		out:       w,
		l:         log.New(w, "", lflag),
		aw:        aw,
		compilers: map[string][]compiler.Named{},
	}
	return l, nil
}

// Close closes the log observer.
func (j *Logger) Close() error {
	return j.out.Close()
}

// OnPrepare logs the preparation attempts of a director.
func (j *Logger) OnPrepare(qs quantity.RootSet, _ pathset.Pathset) {
	qs.Log(j.l)
}

// OnMachines logs a machine block.
func (j *Logger) OnMachines(m machine.Message) {
	switch m.Kind {
	case machine.MessageStart:
		j.l.Printf("%s:\n", stringhelp.PluralQuantity(m.Index, "machine", "", "s"))
	case machine.MessageRecord:
		// TODO(@MattWindsor91): store more information?
		j.l.Printf(" - %s (%s)\n", m.Machine.ID, stringhelp.PluralQuantity(m.Machine.Cores, "core", "", "s"))
	}
}

// OnCycleAnalysis logs s to this logger's file.
func (j *Logger) OnCycleAnalysis(s director.CycleAnalysis) {
	if err := j.aw.WriteSourced(s); err != nil {
		j.l.Println("error writing analysis:", err)
	}
}

// OnCycleSave logs s to this logger's file.
func (j *Logger) OnCycleSave(c director.Cycle, s saver.ArchiveMessage) {
	switch s.Kind {
	case saver.ArchiveStart:
		j.l.Printf("saving (cycle %s) %s to %s\n", c, s.SubjectName, s.File)
	case saver.ArchiveFileMissing:
		j.l.Printf("when saving (cycle %s) %s: missing file %s\n", c, s.SubjectName, s.File)
	}
}

func (j *Logger) OnCycleCompiler(c director.Cycle, m compiler.Message) {
	// TODO(@MattWindsor91): abstract this?
	cs := c.String()
	switch m.Kind {
	case observing.BatchStart:
		j.compilers[cs] = make([]compiler.Named, m.Num)
	case observing.BatchStep:
		j.compilers[cs][m.Num] = *m.Configuration
	case observing.BatchEnd:
		j.logCompilers(c, j.compilers[cs])
	}
}

// logCompilers logs compilers to this Logger's file.
func (j *Logger) logCompilers(c director.Cycle, cs []compiler.Named) {
	j.l.Printf("%s compilers %d:\n", c, len(cs))
	for _, c := range cs {
		j.l.Printf("- %s: %s\n", c.ID, c.Configuration)
	}
}
