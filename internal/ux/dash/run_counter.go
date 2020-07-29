// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/helper/errhelp"

	"github.com/MattWindsor91/act-tester/internal/model/run"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/widgets/sparkline"
	"github.com/mum4k/termdash/widgets/text"
)

// runCounter holds the current run information, as well as a sparkline of previous run times.
type runCounter struct {
	last run.Run

	text  *text.Text
	spark *sparkline.SparkLine
}

func newRunCounter() (*runCounter, error) {
	var (
		r   runCounter
		err error
	)

	if r.text, err = text.New(); err != nil {
		return nil, err
	}
	if r.spark, err = sparkline.New(sparkline.Label("last times")); err != nil {
		return nil, err
	}
	return &r, err
}

// grid outputs a grid arrangement containing this run counter's widgets.
func (r *runCounter) grid() []grid.Element {
	return []grid.Element{
		grid.RowHeightFixed(1, grid.Widget(r.text)),
		grid.RowHeightFixed(2, grid.Widget(r.spark)),
	}
}

func (r *runCounter) onIteration(run run.Run) error {
	err := errhelp.FirstError(r.updateText(run), r.updateSpark(run))
	r.last = run
	return err
}

func (r *runCounter) updateText(run run.Run) error {
	txt := fmt.Sprintf("#%d %s", run.Iter, run.Start.Format(time.Stamp))
	return r.text.Write(txt, text.WriteReplace())
}

func (r *runCounter) updateSpark(run run.Run) error {
	if r.last.Start.IsZero() {
		return nil
	}
	return r.spark.Add([]int{int(run.Start.Sub(r.last.Start).Seconds())})
}
