// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"
	"time"

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

	for _, b := range []struct {
		name   string
		bucket corpus.Corpus
	}{
		{name: "compile failures", bucket: c.CompileFailures},
		{name: "run failures", bucket: c.RunFailures},
		{name: "timeouts", bucket: c.Timeouts},
	} {
		if err := r.logBucket(b.name, b.bucket); err != nil {
			return err
		}
	}
	return nil
}

func (r *ResultLog) logBucket(name string, bucket corpus.Corpus) error {
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
	ncfails := len(c.CompileFailures)
	ntimeouts := len(c.Timeouts)
	nrfails := len(c.RunFailures)
	nsuccesses := len(c.Successes)

	return r.log.Write(
		fmt.Sprintf(
			"[%s #%d %s] %d success, %d c/fail, %d t/out, %d r/fail\n",
			mid.String(),
			iter,
			start.Format(time.Stamp),
			nsuccesses,
			ncfails,
			ntimeouts,
			nrfails,
		),
	)
}
