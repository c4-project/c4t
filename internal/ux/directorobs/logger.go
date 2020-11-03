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

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/plan/analysis"
	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"
)

// TODO(@MattWindsor91): merge this with the singleobs logger?

// Logger is a director observer that emits logs to a writer when cycles finish up.
type Logger struct {
	// done is the channel used to signal that a logger has finished.
	done chan struct{}
	// out is the writer to use for logging analyses.
	out io.WriteCloser
	// l is the intermediate logger that sits atop out.
	l *log.Logger
	// aw is the analyser writer used for outputting sourced analyses.
	aw *pretty.Printer
	// fwd receives forwarded observations from instance loggers.
	fwd *ForwardReceiver
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
		done: make(chan struct{}),
		out:  w,
		l:    log.New(w, "", lflag),
		aw:   aw,
	}
	// TODO(@MattWindsor91): plumb in a capacity somehow
	l.fwd = NewForwardReceiver(l.runStep, 0)
	return l, nil
}

// Run runs the log observer.
// It is not re-entrant.
func (j *Logger) Run(ctx context.Context) error {
	defer close(j.done)
	return j.fwd.Run(ctx)
}

// Close closes the log observer.
func (j *Logger) Close() error {
	return j.out.Close()
}

func (j *Logger) runStep(f Forward) error {
	switch f.Kind {
	case ForwardAnalysis:
		return j.logAnalysis(director.CycleAnalysis{
			Cycle:    f.Cycle.Cycle,
			Analysis: *f.Analysis,
		})
	case ForwardCompiler:
		return j.logCompilerMessage(f.Cycle.Cycle, *f.Compiler)
	case ForwardCycle:
		// do nothing for now
	case ForwardSave:
		return j.logSaving(f.Cycle.Cycle, *f.Save)
	}
	return nil
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

// Instance creates an instance logger.
func (j *Logger) Instance(id.ID) (director.InstanceObserver, error) {
	ch := make(chan Forward)
	j.fwd.Add(ch)
	return &InstanceLogger{
		done: j.done,
		fwd:  ch,
	}, nil
}

// logAnalysis logs s to this logger's file.
func (j *Logger) logAnalysis(s director.CycleAnalysis) error {
	return j.aw.WriteSourced(s)
}

// logSaving logs s to this logger's file.
func (j *Logger) logSaving(c director.Cycle, s saver.ArchiveMessage) error {
	switch s.Kind {
	case saver.ArchiveStart:
		j.l.Printf("saving (cycle %s) %s to %s\n", c, s.SubjectName, s.File)
	case saver.ArchiveFileMissing:
		j.l.Printf("when saving (cycle %s) %s: missing file %s\n", c, s.SubjectName, s.File)
	}
	return nil
}

func (j *Logger) logCompilerMessage(c director.Cycle, m compiler.Message) error {
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
	return nil
}

// logCompilers logs compilers to this Logger's file.
func (j *Logger) logCompilers(c director.Cycle, cs []compiler.Named) {
	j.l.Printf("%s compilers %d:\n", c, len(cs))
	for _, c := range cs {
		j.l.Printf("- %s: %s\n", c.ID, c.Configuration)
	}
}

// InstanceLogger holds state for logging a particular instance.
type InstanceLogger struct {
	// done is a channel closed when the instance can no longer log.
	done <-chan struct{}
	// fwd is a channel to which the instance forwards messages to the main logger.
	fwd chan<- Forward
	// cycle contains information about the current iteration.
	cycle director.Cycle
}

func (l *InstanceLogger) OnCompilerConfig(m compiler.Message) {
	l.forward(ForwardCompilerMessage(l.cycle, m))
}

// OnIteration notes that the instance's iteration has changed.
func (l *InstanceLogger) OnCycle(r director.CycleMessage) {
	if r.Kind == director.CycleStart {
		l.cycle = r.Cycle
	}
	l.forward(ForwardCycleMessage(r))
}

// OnCollation logs a collation to this logger.
func (l *InstanceLogger) OnAnalysis(c analysis.Analysis) {
	l.forward(ForwardAnalysisMessage(l.cycle, c))
}

func (l *InstanceLogger) OnArchive(s saver.ArchiveMessage) {
	l.forward(ForwardSaveMessage(l.cycle, s))
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

// OnInstanceClose doesn't (yet) log the instance closure, but closes the forwarding channel.
func (l *InstanceLogger) OnInstanceClose() {
	// TODO(@MattWindsor91): forward a message about this?
	close(l.fwd)
}

func (l *InstanceLogger) forward(f Forward) {
	select {
	case <-l.done:
	case l.fwd <- f:
	}
}
