// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"github.com/MattWindsor91/act-tester/internal/model/run"
	"github.com/MattWindsor91/act-tester/internal/model/status"
	"github.com/MattWindsor91/act-tester/internal/stage/analyser/observer"

	"github.com/MattWindsor91/act-tester/internal/plan/analyser"

	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/mum4k/termdash/widgets/text"
)

// Observer represents a single machine instance inside a dash.
type Observer struct {
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
func NewObserver(rlog *ResultLog) (*Observer, error) {
	var err error

	d := Observer{
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

// AddToGrid adds this observer to a grid builder gb with the container ID id..
func (o *Observer) AddToGrid(gb *grid.Builder, id string, pc int) {
	gb.Add(grid.RowHeightPercWithOpts(pc,
		[]container.Option{
			container.ID(id),
			container.Border(linestyle.Double),
		},
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

// OnAnalysis observes an analyser by adding failure/timeout/flag rates to the sparklines.
func (o *Observer) OnAnalysis(a analyser.Analysis) {
	for i := status.Ok; i <= status.Last; i++ {
		o.sendStatusCount(i, len(a.ByStatus[i]))
	}
	if err := o.logAnalysis(a); err != nil {
		o.logError(err)
	}
}

// OnArchive currently ignores a save observation.
func (o *Observer) OnArchive(observer.ArchiveMessage) {
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

func (o *Observer) logAnalysis(a analyser.Analysis) error {
	sc := analyser.AnalysisWithRun{
		Run:      o.run.last,
		Analysis: a,
	}
	return o.rlog.Log(sc)
}

// OnBuild forwards a build observation.
func (o *Observer) OnBuild(m builder.Message) {
	switch m.Kind {
	case builder.BuildStart:
		o.action.OnBuildStart(*m.Manifest)
	case builder.BuildRequest:
		o.action.OnBuildRequest(*m.Request)
	case builder.BuildFinish:
		o.action.OnBuildFinish()
	}
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
