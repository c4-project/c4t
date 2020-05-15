// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package observer

import (
	"context"
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/controller/query"

	"github.com/MattWindsor91/act-tester/internal/model/run"

	"github.com/MattWindsor91/act-tester/internal/model/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"
)

// BasicLogger is an interface for things that can use the Log func.
type BasicLogger interface {
	// LogHeader should log the String() of the argument collation.
	LogHeader(analysis.Sourced) error

	// LogBucketHeader should log a header for the collation bucket with the given status.
	LogBucketHeader(subject.Status) error

	// LogBucketEntry should log the given subject name, assuming LogBucketHeader has been called.
	LogBucketEntry(string) error
}

// Log logs s to b.
func Log(b BasicLogger, s analysis.Sourced) error {
	if err := b.LogHeader(s); err != nil {
		return err
	}
	return logBuckets(b, s)
}

func logBuckets(b BasicLogger, s analysis.Sourced) error {
	sc := s.Collation.ByStatus
	for i := subject.FirstBadStatus; i < subject.NumStatus; i++ {
		if err := logBucket(b, i, sc[i]); err != nil {
			return err
		}
	}
	return nil
}

func logBucket(b BasicLogger, s subject.Status, bucket corpus.Corpus) error {
	if len(bucket) == 0 {
		return nil
	}
	if err := b.LogBucketHeader(s); err != nil {
		return err
	}
	for _, n := range bucket.Names() {
		if err := b.LogBucketEntry(n); err != nil {
			return err
		}
	}
	return nil
}

// Logger is a director observer that emits logs to a writer when runs finish up.
type Logger struct {
	// out is the writer to use for logging collations.
	out io.Writer
	// aw is the analysis writer used for outputting sourced analyses.
	aw *query.AnalysisWriter
	// collCh is used to send sourced analyses for logging.
	collCh chan analysis.Sourced
	// compCh is used to send compilers for logging.
	compCh chan compilerSet
}

// NewLogger constructs a new Logger writing into w, ranging over machine IDs ids.
func NewLogger(w io.Writer) (*Logger, error) {
	aw, err := query.NewAnalysisWriter(&query.Config{Out: w})
	if err != nil {
		return nil, err
	}
	return &Logger{out: w, aw: aw, collCh: make(chan analysis.Sourced), compCh: make(chan compilerSet)}, nil
}

// Run runs the log observer.
func (j *Logger) Run(ctx context.Context, _ func()) error {
	for {
		if err := j.runStep(ctx); err != nil {
			return err
		}
	}
}

func (j *Logger) runStep(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case sc := <-j.collCh:
		return j.logCollation(sc)
	case cc := <-j.compCh:
		return j.logCompilers(cc)
	}
}

// Instance creates an instance logger.
func (j *Logger) Instance(id.ID) (Instance, error) {
	return &InstanceLogger{collCh: j.collCh, compCh: j.compCh}, nil
}

// logCollation logs s to this Logger's file.
func (j *Logger) logCollation(s analysis.Sourced) error {
	return j.aw.WriteSourced(&s)
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
	// collCh is the channel used to send sourced collations for logging.
	collCh chan<- analysis.Sourced
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
func (l *InstanceLogger) OnCollation(c *analysis.Analysis) {
	select {
	case <-l.done:
	case l.collCh <- l.addSource(c):
	}
}

func (l *InstanceLogger) addSource(c *analysis.Analysis) analysis.Sourced {
	return analysis.Sourced{
		Run:       l.run,
		Collation: c,
	}
}

// OnBuildStart does nothing.
func (l *InstanceLogger) OnBuildStart(builder.Manifest) {}

// OnBuildRequest does nothing.
func (l *InstanceLogger) OnBuildRequest(builder.Request) {}

// OnBuildFinish does nothing.
func (l *InstanceLogger) OnBuildFinish() {}

// OnCopyStart does nothing.
func (l *InstanceLogger) OnCopyStart(int) {}

// OnCopy does nothing.
func (l *InstanceLogger) OnCopy(string, string) {}

// OnCopyFinish does nothing.
func (l *InstanceLogger) OnCopyFinish() {}
