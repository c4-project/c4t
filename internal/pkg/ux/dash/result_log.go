// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/collate"
	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"
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

// Log logs a collation for the machine with ID mid on iteration iter, starting at start.
func (r *ResultLog) Log(mid id.ID, iter uint64, start time.Time, c *collate.Collation) error {
	if err := r.logHeader(mid, iter, start, c); err != nil {
		return err
	}

	sc := c.ByStatus()
	for i := subject.FirstBadStatus; i < subject.NumStatus; i++ {
		if err := r.logBucket(i.String(), sc[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *ResultLog) logBucket(name string, bucket corpus.Corpus) error {
	if len(bucket) == 0 {
		return nil
	}
	header := fmt.Sprintf("  [%s]\n", name)
	if err := r.log.Write(header, text.WriteCellOpts(cell.FgColor(colorFailed))); err != nil {
		return err
	}
	for _, n := range bucket.Names() {
		if err := r.log.Write(fmt.Sprintf("  - %s\n", n)); err != nil {
			return err
		}
	}
	return nil
}

func (r *ResultLog) logHeader(mid id.ID, iter uint64, start time.Time, c *collate.Collation) error {
	sc := collate.Sourced{
		MachineID: mid,
		Iter:      iter,
		Start:     start,
		Collation: c,
	}
	return r.log.Write(sc.String())
}
