// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package observer

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/collate"
	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// BasicLogger is an interface for things that can use the LogCollation func.
type BasicLogger interface {
	// LogHeader should log the String() of the argument collation.
	LogHeader(collate.Sourced) error

	// LogBucketHeader should log a header for the collation bucket with the given status.
	LogBucketHeader(subject.Status) error

	// LogBucketEntry should log the given subject name, assuming LogBucketHeader has been called.
	LogBucketEntry(string) error
}

// Writer wraps io.Writer to implement BasicLogger.
type Writer struct {
	io.Writer
}

// LogHeader logs a collation header to the writer.
func (w Writer) LogHeader(s collate.Sourced) error {
	_, err := fmt.Fprintln(w.Writer, &s)
	return err
}

// LogBucketHeader logs a collation bucket header to the writer.
func (w Writer) LogBucketHeader(s subject.Status) error {
	_, err := fmt.Fprintf(w.Writer, "  [%s]\n", s)
	return err
}

// LogBucketEntry logs a collation subject name to the writer.
func (w Writer) LogBucketEntry(sname string) error {
	_, err := fmt.Fprintf(w.Writer, "  - %s\n", sname)
	return err
}

// log logs s to b.
func Log(b BasicLogger, s collate.Sourced) error {
	if err := b.LogHeader(s); err != nil {
		return err
	}
	return logBuckets(b, s)
}

func logBuckets(b BasicLogger, s collate.Sourced) error {
	sc := s.Collation.ByStatus()
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
	// ch is used to send SourcedCollations for logging.
	ch chan collate.Sourced
}

// NewLogger constructs a new LogObserver writing into w, ranging over machine IDs ids.
func NewLogger(w io.Writer) *Logger {
	return &Logger{out: w, ch: make(chan collate.Sourced)}
}

// Run runs the log observer.
func (j *Logger) Run(ctx context.Context, _ func()) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sc := <-j.ch:
			if err := j.log(sc); err != nil {
				return err
			}
		}
	}
}

// Instance creates an instance logger for machine mid.
func (j *Logger) Instance(mid id.ID) (Instance, error) {
	return &InstanceLogger{id: mid, ch: j.ch}, nil
}

// log logs s to this Logger's file.
func (j *Logger) log(s collate.Sourced) error {
	return Log(Writer{j.out}, s)
}

// InstanceLogger holds state for logging a particular instance.
type InstanceLogger struct {
	// done is a channel closed when the instance can no longer log.
	done <-chan struct{}
	// ch is the channel used to send sourced collations for logging.
	ch chan<- collate.Sourced
	// id is the machine ID.
	id id.ID
	// iter is the number of the current iteration.
	iter uint64
	// start is the start time of the current iteration.
	start time.Time
}

// OnIteration notes that the instance's iteration has changed.
func (i *InstanceLogger) OnIteration(iter uint64, start time.Time) {
	i.iter = iter
	i.start = start
}

// OnCollation logs a collation to this logger.
func (i *InstanceLogger) OnCollation(c *collate.Collation) {
	select {
	case <-i.done:
	case i.ch <- i.addSource(c):
	}
}

func (i *InstanceLogger) addSource(c *collate.Collation) collate.Sourced {
	return collate.Sourced{
		MachineID: i.id,
		Iter:      i.iter,
		Start:     i.start,
		Collation: c,
	}
}

// OnBuildStart does nothing.
func (i *InstanceLogger) OnBuildStart(builder.Manifest) {}

// OnBuildRequest does nothing.
func (i *InstanceLogger) OnBuildRequest(builder.Request) {}

// OnBuildFinish does nothing.
func (i *InstanceLogger) OnBuildFinish() {}

// OnCopyStart does nothing.
func (i *InstanceLogger) OnCopyStart(int) {}

// OnCopy does nothing.
func (i *InstanceLogger) OnCopy(string, string) {}

// OnCopyFinish does nothing.
func (i *InstanceLogger) OnCopyFinish() {}
