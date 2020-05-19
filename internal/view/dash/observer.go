// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"github.com/MattWindsor91/act-tester/internal/controller/analyse/observer"
	"github.com/MattWindsor91/act-tester/internal/model/run"
	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"

	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/mum4k/termdash/widgets/text"
)

// Observer represents a single machine instance inside a dash.
type Observer struct {
	mid id.ID

	run   *runCounter
	rlog  *ResultLog
	tally *tally

	sparks *sparkset

	action *actionObserver

	// compilers contains a readout of the currently planned compilers for this instance.
	compilers *text.Text

	nruns uint64
}

// NewObserver constructs an Observer, initialising its various widgets.
func NewObserver(mid id.ID, rlog *ResultLog) (*Observer, error) {
	var err error

	d := Observer{
		mid:  mid,
		rlog: rlog,
	}

	if d.tally, err = newTally(); err != nil {
		return nil, err
	}

	if d.run, err = newRunCounter(); err != nil {
		return nil, err
	}

	if d.sparks, err = newSparkset(); err != nil {
		return nil, err
	}

	if d.action, err = NewCorpusObserver(); err != nil {
		return nil, err
	}

	if d.compilers, err = text.New(text.WrapAtWords()); err != nil {
		return nil, err
	}

	return &d, nil
}

const (
	percRun     = 25
	percStats   = 25
	percActions = 100 - percRun - percStats
)

// AddToGrid adds this observer to a grid builder gb.
func (o *Observer) AddToGrid(gb *grid.Builder, midstr string, pc int) {
	gb.Add(grid.RowHeightPercWithOpts(pc,
		[]container.Option{container.Border(linestyle.Double), container.BorderTitle(midstr)},
		grid.ColWidthPerc(percRun,
			grid.RowHeightPercWithOpts(
				40,
				[]container.Option{
					container.Border(linestyle.Light),
					container.BorderTitle("Run"),
				},
				o.run.grid()...,
			),
			grid.RowHeightPerc(60, grid.Widget(o.compilers, container.Border(linestyle.Light), container.BorderTitle("Compilers"))),
		),
		grid.ColWidthPerc(percStats,
			grid.RowHeightPercWithOpts(
				40,
				[]container.Option{container.Border(linestyle.Light), container.BorderTitle("Statistics")},
				o.tally.grid()...,
			),
			grid.RowHeightPercWithOpts(60,
				[]container.Option{container.Border(linestyle.Light), container.BorderTitle("Sparklines")},
				o.sparks.grid()...),
		),
		o.currentRunColumn(),
	))
}

func (o *Observer) currentRunColumn() grid.Element {
	return grid.ColWidthPercWithOpts(percActions,
		[]container.Option{
			container.Border(linestyle.Light),
			container.BorderTitle("Current Run"),
		},
		o.action.gridRows()...,
	)
}

// OnIteration logs that a new iteration has begun.
func (o *Observer) OnIteration(r run.Run) {
	o.nruns = r.Iter
	_ = o.run.onIteration(r)
	o.action.reset()
}

// OnAnalysis observes an analysis by adding failure/timeout/flag rates to the sparklines.
func (o *Observer) OnAnalysis(a analysis.Analysis) {
	for i := status.Ok; i < status.Num; i++ {
		o.sendStatusCount(i, len(a.ByStatus[i]))
	}
	if err := o.logAnalysis(a); err != nil {
		o.logError(err)
	}
}

// OnSave currently ignores a save observation.
func (o *Observer) OnSave(observer.Saving) {
	// TODO(@MattWindsor91): do something with this?
}

// OnSaveFileMissing currently ignores a save missing-file observation.
func (o *Observer) OnSaveFileMissing(observer.Saving, string) {
	// TODO(@MattWindsor91): do something with this?
}

func (o *Observer) sendStatusCount(i status.Status, n int) {
	if err := o.tally.tallyStatus(i, n); err != nil {
		o.logError(err)
	}
	if err := o.sparks.sparkStatus(i, n); err != nil {
		o.logError(err)
	}
}

func (o *Observer) logAnalysis(a analysis.Analysis) error {
	sc := analysis.Sourced{
		Run:      o.run.last,
		Analysis: a,
	}
	return o.rlog.Log(sc)
}

// OnBuildStart forwards a build start observation.
func (o *Observer) OnBuildStart(m builder.Manifest) {
	o.action.OnBuildStart(m)
}

// OnBuildRequest forwards a build request observation.
func (o *Observer) OnBuildRequest(r builder.Request) {
	o.action.OnBuildRequest(r)
}

// OnBuildFinish forwards a build finish observation.
func (o *Observer) OnBuildFinish() {
	o.action.OnBuildFinish()
}

// OnCopyStart forwards a copy start observation.
func (o *Observer) OnCopyStart(nfiles int) {
	o.action.OnCopyStart(nfiles)
}

// OnCopy forwards a copy observation.
func (o *Observer) OnCopy(dst, src string) {
	o.action.OnCopy(dst, src)
}

// OnCopyFinish forwards a copy finish observation.
func (o *Observer) OnCopyFinish() {
	o.action.OnCopyFinish()
}

func (o *Observer) logError(err error) {
	// For want of better location.
	o.action.logError(err)
}
