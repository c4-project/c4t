// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"github.com/c4-project/c4t/internal/subject/status"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/widgets/sparkline"
)

// sparkset contains the sparklines for a machine.
type sparkset struct {
	// statusLines contains one sparkline for each status.
	// (This _includes_ StatusUnknown to simplify calculations later, but we don't display it as a line.)
	statusLines [status.Last + 1]*sparkline.SparkLine
}

func newSparkset() (*sparkset, error) {
	var (
		s   sparkset
		err error
	)

	for i := status.Ok; i <= status.Last; i++ {
		if s.statusLines[i], err = sparkline.New(
			sparkline.Color(statusColours[i]), sparkline.Label(i.String()),
		); err != nil {
			return nil, err
		}
	}
	return &s, err
}

func (s *sparkset) sparkLines() []*sparkline.SparkLine {
	return s.statusLines[status.Ok:]
}

func (s *sparkset) grid() []grid.Element {
	sls := s.sparkLines()
	els := make([]grid.Element, len(sls))
	for i, sl := range sls {
		els[i] = grid.RowHeightFixed(2, grid.Widget(sl))
	}
	return els
}

func (s *sparkset) sparkStatus(st status.Status, n int) error {
	return s.statusLines[st].Add([]int{n})
}
