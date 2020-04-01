// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/collate"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/widgets/sparkline"
)

// sparks contains the sparklines for a machine.
type sparkset struct {
	runTime  *sparkline.SparkLine
	cfails   *sparkline.SparkLine
	timeouts *sparkline.SparkLine
	flags    *sparkline.SparkLine
}

func newSparkset() (*sparkset, error) {
	var (
		s   sparkset
		err error
	)
	for _, pp := range []struct {
		p **sparkline.SparkLine
		c cell.Color
		l string
	}{
		{l: "Time ", p: &s.runTime, c: cell.ColorGreen},
		{l: "CFail", p: &s.cfails, c: cell.ColorRed},
		{l: "T/Out", p: &s.timeouts, c: cell.ColorMagenta},
		{l: "Flags", p: &s.flags, c: cell.ColorYellow},
	} {
		if *pp.p, err = sparkline.New(sparkline.Color(pp.c), sparkline.Label(pp.l)); err != nil {
			return nil, err
		}
	}
	return &s, err
}

func (s *sparkset) sparkLines() []*sparkline.SparkLine {
	return []*sparkline.SparkLine{s.runTime, s.cfails, s.timeouts, s.flags}
}

func (s *sparkset) gridRows() []grid.Element {
	sls := s.sparkLines()
	els := make([]grid.Element, len(sls))
	for i, sl := range sls {
		els[i] = grid.RowHeightFixed(2, grid.Widget(sl))
	}
	return els
}

func (s *sparkset) sparkCollation(c *collate.Collation) error {
	ferr := s.flags.Add([]int{len(c.Flagged)})
	terr := s.timeouts.Add([]int{len(c.Run.Timeouts)})
	cerr := s.cfails.Add([]int{len(c.Compile.Failures)})
	return iohelp.FirstError(ferr, terr, cerr)
}
