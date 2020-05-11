// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/director/observer"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/analysis"
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
func (r *ResultLog) Log(sc analysis.Sourced) error {
	return observer.Log(r, sc)
}

// LogHeader logs the header of a sourced collation sc.
func (r *ResultLog) LogHeader(sc analysis.Sourced) error {
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
func (r *ResultLog) LogBucketHeader(st subject.Status) error {
	header := fmt.Sprintf("  [%s]\n", st)
	return r.log.Write(header, text.WriteCellOpts(cell.FgColor(statusColours[st])))
}

// LogHeader logs an entry for a subject with name sname.
func (r *ResultLog) LogBucketEntry(sname string) error {
	return r.log.Write(fmt.Sprintf("  - %s\n", sname))
}
