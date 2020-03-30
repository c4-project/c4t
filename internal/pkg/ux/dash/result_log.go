// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/director/observer"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/collate"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
)

// ResultLog provides a logging mechanism for collated subjects.
type ResultLog struct {
	log *text.Text
}

// NewResultLog creates a new ResultLog.
func NewResultLog() (*ResultLog, error) {
	var (
		r   ResultLog
		err error
	)
	r.log, err = text.New(text.RollContent())
	return &r, err
}

// Log logs a sourced collation sc.
func (r *ResultLog) Log(sc collate.Sourced) error {
	return observer.Log(r, sc)
}

// LogHeader logs the header of a sourced collation sc.
func (r *ResultLog) LogHeader(sc collate.Sourced) error {
	return r.log.Write(sc.String()+"\n", text.WriteCellOpts(cell.FgColor(summaryColor(sc))))
}

// summaryColor retrieves a colour to use for the log header of sc, according to a 'traffic lights' system.
func summaryColor(sc collate.Sourced) cell.Color {
	switch {
	case sc.Collation.HasFailures():
		return colorFailed
	case sc.Collation.HasFlagged():
		return colorFlagged
	default:
		return colorRun
	}
}

// LogHeader logs the header of a collation bucket with status st.
func (r *ResultLog) LogBucketHeader(st subject.Status) error {
	header := fmt.Sprintf("  [%s]\n", st)
	return r.log.Write(header, text.WriteCellOpts(cell.FgColor(colorFailed)))
}

// LogHeader logs an entry for a subject with name sname.
func (r *ResultLog) LogBucketEntry(sname string) error {
	return r.log.Write(fmt.Sprintf("  - %s\n", sname))
}
