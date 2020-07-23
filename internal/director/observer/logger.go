// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package observer

import (
	"context"
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/stage/analyse/pretty"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/machine"

	"github.com/MattWindsor91/act-tester/internal/stage/analyse/observer"

	"github.com/MattWindsor91/act-tester/internal/model/run"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/plan/analysis"
)

// Logger is a director observer that emits logs to a writer when runs finish up.
type Logger struct {
	// out is the writer to use for logging analyses.
	// It will be closed by the director.
	out io.WriteCloser
	// aw is the analysis writer used for outputting sourced analyses.
	aw *pretty.Printer
	// anaCh is used to send sourced analyses for logging.
	anaCh chan analysis.Sourced
	// compCh is used to send compilers for logging.
	compCh chan compilerSet
	// saveCh is used to send save actions for logging.
	saveCh chan archiveMessage
}

// NewLogger constructs a new Logger writing into w, ranging over machine IDs ids.
func NewLogger(w io.WriteCloser) (*Logger, error) {
	aw, err := pretty.NewPrinter(
		pretty.WriteTo(w),
		pretty.ShowCompilers(true),
		pretty.ShowSubjects(true),
	)
	if err != nil {
		return nil, err
	}
	return &Logger{
		out:    w,
		aw:     aw,
		anaCh:  make(chan analysis.Sourced),
		compCh: make(chan compilerSet),
		saveCh: make(chan archiveMessage),
	}, nil
}

// Run runs the log observer.
func (j *Logger) Run(ctx context.Context, _ func()) error {
	for {
		if err := j.runStep(ctx); err != nil {
			return err
		}
	}
}

func (j *Logger) Close() error {
	return j.out.Close()
}

func (j *Logger) runStep(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ac := <-j.anaCh:
		return j.logAnalysis(ac)
	case cc := <-j.compCh:
		return j.logCompilers(cc)
	case sc := <-j.saveCh:
		return j.logSaving(sc)
	}
}

// OnMachines logs a machine block.
func (j *Logger) OnMachines(m machine.Message) {
	switch m.Kind {
	case machine.MessageStart:
		_, _ = fmt.Fprintf(j.out, "%s:\n", iohelp.PluralQuantity(m.Index, "machine", "", "s"))
	case machine.MessageRecord:
		// TODO(@MattWindsor91): store more information?
		_, _ = fmt.Fprintf(j.out, " - %s (%s)\n", m.Machine.ID, iohelp.PluralQuantity(m.Machine.Cores, "core", "", "s"))
	}
}

// Instance creates an instance logger.
func (j *Logger) Instance(id.ID) (Instance, error) {
	return &InstanceLogger{anaCh: j.anaCh, compCh: j.compCh, saveCh: j.saveCh}, nil
}

// logAnalysis logs s to this logger's file.
func (j *Logger) logAnalysis(s analysis.Sourced) error {
	return j.aw.WriteSourced(s)
}

// logSaving logs s to this logger's file.
func (j *Logger) logSaving(s archiveMessage) error {
	var err error
	switch s.body.Kind {
	case observer.ArchiveStart:
		_, err = fmt.Fprintf(j.out, "saving (run %s) %s to %s\n", s.run, s.body.SubjectName, s.body.File)
	case observer.ArchiveFileMissing:
		_, err = fmt.Fprintf(j.out, "when saving (run %s) %s: missing file %s\n", s.run, s.body.SubjectName, s.body.File)
	}
	return err
}

// logCompilers logs compilers to this Logger's file.
func (j *Logger) logCompilers(cs compilerSet) error {
	// TODO(@MattWindsor91): abstract this?
	if _, err := fmt.Fprintf(j.out, "%s compilers %d:\n", cs.run, len(cs.compilers)); err != nil {
		return err
	}
	for _, c := range cs.compilers {
		if _, err := fmt.Fprintf(j.out, "- %s: %s\n", c.ID, c.Compiler); err != nil {
			return err
		}
	}
	return nil
}

// InstanceLogger holds state for logging a particular instance.
type InstanceLogger struct {
	// done is a channel closed when the instance can no longer log.
	done <-chan struct{}
	// compCh is the channel used to send compiler sets for logging.
	compCh chan<- compilerSet
	// anaCh is the channel used to send sourced analyses for logging.
	anaCh chan<- analysis.Sourced
	// saveCh is the channel used to send save actions for logging.
	saveCh chan<- archiveMessage
	// run contains information about the current iteration.
	run run.Run
	// compilers stores the current, if any, compiler set.
	compilers []compiler.Named
	// icompiler stores the index of the compiler being received.
	icompiler int
}

type compilerSet struct {
	run       run.Run
	compilers []compiler.Named
}

type archiveMessage struct {
	run  run.Run
	body observer.ArchiveMessage
}

func (l *InstanceLogger) OnCompilerPlanStart(ncompilers int) {
	l.compilers = make([]compiler.Named, ncompilers)
	l.icompiler = 0
}

func (l *InstanceLogger) OnCompilerPlan(c compiler.Named) {
	l.compilers[l.icompiler] = c
	l.icompiler++
}

func (l *InstanceLogger) OnCompilerPlanFinish() {
	select {
	case <-l.done:
	case l.compCh <- l.makeCompilerSet():
	}
}

func (l *InstanceLogger) makeCompilerSet() compilerSet {
	return compilerSet{
		run:       l.run,
		compilers: l.compilers,
	}
}

// OnIteration notes that the instance's iteration has changed.
func (l *InstanceLogger) OnIteration(r run.Run) {
	l.run = r
}

// OnCollation logs a collation to this logger.
func (l *InstanceLogger) OnAnalysis(c analysis.Analysis) {
	select {
	case <-l.done:
	case l.anaCh <- l.addSource(c):
	}
}

func (l *InstanceLogger) OnArchive(s observer.ArchiveMessage) {
	msg := archiveMessage{
		run:  l.run,
		body: s,
	}
	select {
	case <-l.done:
	case l.saveCh <- msg:
	}
}

func (l *InstanceLogger) addSource(c analysis.Analysis) analysis.Sourced {
	return analysis.Sourced{
		Run:      l.run,
		Analysis: c,
	}
}

// OnBuild does nothing.
func (l *InstanceLogger) OnBuild(builder.Message) {}

// OnCopyStart does nothing.
func (l *InstanceLogger) OnCopyStart(int) {}

// OnCopy does nothing.
func (l *InstanceLogger) OnCopy(string, string) {}

// OnCopyFinish does nothing.
func (l *InstanceLogger) OnCopyFinish() {}
