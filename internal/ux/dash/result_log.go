// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"
	"time"

	"github.com/c4-project/c4t/internal/director"

	"github.com/c4-project/c4t/internal/subject/corpus"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
)

// resultLogOverflow is the number of written headers at which the result log clears itself (to avoid eating memory).
const resultLogOverflow uint16 = 1000

// ResultLog provides a logging mechanism for collated subjects.
type ResultLog struct {
	log      *text.Text
	nheaders uint16
}

// NewResultLog creates a new ResultLog.
func NewResultLog() (*ResultLog, error) {
	var (
		r   ResultLog
		err error
	)
	r.log, err = text.New(text.RollContent(), text.WrapAtWords())
	return &r, err
}

// Log logs a sourced collation sc.
func (r *ResultLog) Log(sc director.CycleAnalysis) error {
	if err := r.LogHeader(sc); err != nil {
		return err
	}
	return r.logBuckets(sc)
}

func (r *ResultLog) logBuckets(s director.CycleAnalysis) error {
	sc := s.Analysis.ByStatus
	for i := status.FirstBad; i <= status.Last; i++ {
		if err := r.logBucket(i, sc[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *ResultLog) logBucket(s status.Status, bucket corpus.Corpus) error {
	if len(bucket) == 0 {
		return nil
	}
	if err := r.LogBucketHeader(s); err != nil {
		return err
	}
	for _, n := range bucket.Names() {
		if err := r.LogBucketEntry(n); err != nil {
			return err
		}
	}
	return nil
}

// LogHeader logs the header of a sourced collation sc.
func (r *ResultLog) LogHeader(sc director.CycleAnalysis) error {
	if err := r.maybeOverflow(); err != nil {
		return err
	}
	return r.log.Write(sc.String()+"\n", text.WriteCellOpts(cell.FgColor(summaryColor(sc))))
}

func (r *ResultLog) maybeOverflow() error {
	r.nheaders++
	if r.nheaders < resultLogOverflow {
		return nil
	}
	r.nheaders = 0
	return r.overflow()
}

func (r *ResultLog) overflow() error {
	txt := fmt.Sprintf("[log overflowed at %s; see log file]\n", time.Now().Format(time.Stamp))
	return r.log.Write(txt, text.WriteReplace(), text.WriteCellOpts(cell.FgColor(cell.ColorMagenta)))
}

// LogHeader logs the header of a collation bucket with status st.
func (r *ResultLog) LogBucketHeader(st status.Status) error {
	header := fmt.Sprintf("  [%s]\n", st)
	return r.log.Write(header, text.WriteCellOpts(cell.FgColor(statusColours[st])))
}

// LogHeader logs an entry for a subject with name sname.
func (r *ResultLog) LogBucketEntry(sname string) error {
	return r.log.Write(fmt.Sprintf("  - %s\n", sname))
}
