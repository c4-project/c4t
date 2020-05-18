// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/widgets/text"
)

type tally struct {
	nstatus [status.Num]uint64
	dstatus [status.Num]*text.Text
}

func newTally() (*tally, error) {
	var (
		t   tally
		err error
	)

	for i := status.Ok; i < status.Num; i++ {
		if t.dstatus[i], err = text.New(text.DisableScrolling()); err != nil {
			return nil, err
		}
	}
	return &t, nil
}

func (t *tally) grid() []grid.Element {
	widgets := t.dstatus[status.Ok:]
	els := make([]grid.Element, len(widgets))
	for i, w := range widgets {
		els[i] = grid.RowHeightFixed(1, grid.Widget(w))
	}
	return els
}

func (t *tally) tallyCollation(c *analysis.Analysis) error {
	for i := status.Ok; i < status.Num; i++ {
		t.nstatus[i] += uint64(len(c.ByStatus[i]))
		if err := t.updateCollation(i); err != nil {
			return err
		}
	}
	return nil
}

func (t *tally) updateCollation(s status.Status) error {
	if err := t.dstatus[s].Write(
		s.String(), text.WriteCellOpts(cell.FgColor(statusColours[s])), text.WriteReplace(),
	); err != nil {
		return err
	}
	return t.dstatus[s].Write(fmt.Sprintf(": %d", t.nstatus[s]))
}
