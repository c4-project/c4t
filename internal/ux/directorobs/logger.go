// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package directorobs

import (
	"context"
	"io"
	"log"

	"github.com/MattWindsor91/act-tester/internal/director/pathset"
	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/director"

	"github.com/MattWindsor91/act-tester/internal/copier"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer"

	"github.com/MattWindsor91/act-tester/internal/helper/stringhelp"

	"github.com/MattWindsor91/act-tester/internal/stage/planner"

	"github.com/MattWindsor91/act-tester/internal/stage/perturber"

	"github.com/MattWindsor91/act-tester/internal/observing"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/pretty"

	"github.com/MattWindsor91/act-tester/internal/machine"

	"github.com/MattWindsor91/act-tester/internal/model/run"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/plan/analysis"
	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"
)

// TODO(@MattWindsor91): merge this with the singleobs logger?

// Logger is a director observer that emits logs to a writer when cycles finish up.
type Logger struct {
	// out is the writer to use for logging analyses.
	out io.WriteCloser
	// l is the intermediate logger that sits atop out.
	l *log.Logger
	// aw is the analyser writer used for outputting sourced analyses.
	aw *pretty.Printer
	// anaCh is used to send sourced analyses for logging.
	anaCh chan analysis.WithRun
	// compCh is used to send compilers for logging.
	compCh chan compilerSet
	// saveCh is used to send save actions for logging.
	saveCh chan archiveMessage
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
	return &Logger{
		out:    w,
		l:      log.New(w, "", lflag),
		aw:     aw,
		anaCh:  make(chan analysis.WithRun),
		compCh: make(chan compilerSet),
		saveCh: make(chan archiveMessage),
	}, nil
}

// Run runs the log observer.
func (j *Logger) Run(ctx context.Context) error {
	for {
		if err := j.runStep(ctx); err != nil {
			return err
		}
	}
}

// Close closes the log observer.
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
		j.logCompilers(cc)
	case sc := <-j.saveCh:
		j.logSaving(sc)
	}
	return nil
}

// OnPrepare logs the preparation attempts of a director.
func (j *Logger) OnPrepare(qs quantity.RootSet, ps pathset.Pathset) {
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

// Instance creates an instance logger.
func (j *Logger) Instance(id.ID) (director.InstanceObserver, error) {
	return &InstanceLogger{anaCh: j.anaCh, compCh: j.compCh, saveCh: j.saveCh}, nil
}

// logAnalysis logs s to this logger's file.
func (j *Logger) logAnalysis(s analysis.WithRun) error {
	return j.aw.WriteSourced(s)
}

// logSaving logs s to this logger's file.
func (j *Logger) logSaving(s archiveMessage) {
	switch s.body.Kind {
	case saver.ArchiveStart:
		j.l.Printf("saving (run %s) %s to %s\n", s.run, s.body.SubjectName, s.body.File)
	case saver.ArchiveFileMissing:
		j.l.Printf("when saving (run %s) %s: missing file %s\n", s.run, s.body.SubjectName, s.body.File)
	}
}

// logCompilers logs compilers to this Logger's file.
func (j *Logger) logCompilers(cs compilerSet) {
	// TODO(@MattWindsor91): abstract this?
	j.l.Printf("%s compilers %d:\n", cs.run, len(cs.compilers))
	for _, c := range cs.compilers {
		j.l.Printf("- %s: %s\n", c.ID, c.Configuration)
	}
}

// InstanceLogger holds state for logging a particular instance.
type InstanceLogger struct {
	// done is a channel closed when the instance can no longer log.
	done <-chan struct{}
	// compCh is the channel used to send compiler sets for logging.
	compCh chan<- compilerSet
	// anaCh is the channel used to send sourced analyses for logging.
	anaCh chan<- analysis.WithRun
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
	body saver.ArchiveMessage
}

func (l *InstanceLogger) OnCompilerConfig(m compiler.Message) {
	switch m.Kind {
	case observing.BatchStart:
		l.onCompilerPlanStart(m.Num)
	case observing.BatchStep:
		l.onCompilerPlan(*m.Configuration)
	case observing.BatchEnd:
		l.onCompilerPlanFinish()
	}
}

func (l *InstanceLogger) onCompilerPlanStart(ncompilers int) {
	l.compilers = make([]compiler.Named, ncompilers)
	l.icompiler = 0
}

func (l *InstanceLogger) onCompilerPlan(c compiler.Named) {
	l.compilers[l.icompiler] = c
	l.icompiler++
}

func (l *InstanceLogger) onCompilerPlanFinish() {
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

func (l *InstanceLogger) OnArchive(s saver.ArchiveMessage) {
	msg := archiveMessage{
		run:  l.run,
		body: s,
	}
	select {
	case <-l.done:
	case l.saveCh <- msg:
	}
}

func (l *InstanceLogger) addSource(c analysis.Analysis) analysis.WithRun {
	return analysis.WithRun{
		Run:      l.run,
		Analysis: c,
	}
}

// OnPerturb does nothing, at the moment.
func (l *InstanceLogger) OnPerturb(perturber.Message) {}

// OnPlan does nothing, at the moment.
func (l *InstanceLogger) OnPlan(planner.Message) {}

// OnBuild does nothing.
func (l *InstanceLogger) OnBuild(builder.Message) {}

// OnCopy does nothing.
func (l *InstanceLogger) OnCopy(copier.Message) {}

// OnMachineNodeMessage does nothing.
func (l *InstanceLogger) OnMachineNodeAction(observer.Message) {}
