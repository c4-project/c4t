// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"strconv"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/run"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/collate"

	"github.com/mum4k/termdash/widgets/segmentdisplay"

	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/mum4k/termdash/widgets/text"
)

// Observer is a BuilderObserver that attaches into a Dash.
type Observer struct {
	mid id.ID

	rlog *ResultLog

	runCount *segmentdisplay.SegmentDisplay
	runStart *text.Text
	sparks   *sparkset

	action *actionObserver

	lastTime time.Time

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

	if d.runCount, err = segmentdisplay.New(); err != nil {
		return nil, err
	}

	if d.runStart, err = text.New(text.DisableScrolling()); err != nil {
		return nil, err
	}

	if d.sparks, err = newSparkset(); err != nil {
		return nil, err
	}

	if d.action, err = NewCorpusObserver(); err != nil {
		return nil, err
	}

	if d.compilers, err = text.New(); err != nil {
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
			grid.RowHeightFixed(10, grid.Widget(o.runCount, container.Border(linestyle.Light), container.BorderTitle("Run#"))),
			grid.RowHeightFixed(1, grid.Widget(o.compilers, container.Border(linestyle.Light), container.BorderTitle("Compilers"))),
		),
		grid.ColWidthPerc(percStats,
			grid.RowHeightFixed(3, grid.Widget(o.runStart, container.Border(linestyle.Light), container.BorderTitle("Start"))),
			grid.RowHeightFixedWithOpts(10,
				[]container.Option{container.Border(linestyle.Light), container.BorderTitle("Sparklines")},
				o.sparks.gridRows()...),
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
	ch := segmentdisplay.NewChunk(strconv.FormatUint(o.nruns, 10))
	_ = o.runCount.Write(
		[]*segmentdisplay.TextChunk{ch},
	)

	_ = o.runStart.Write(r.Start.Format(time.Stamp), text.WriteReplace())

	o.addDurationToSparkline(r.Start)
	o.action.reset()
}

func (o *Observer) addDurationToSparkline(t time.Time) {
	if !o.lastTime.IsZero() {
		dur := t.Sub(o.lastTime)
		_ = o.sparks.runLine.Add([]int{int(dur.Seconds())})
	}
	o.lastTime = t
}

// OnCollation observes a collation by adding failure/timeout/flag rates to the sparklines.
func (o *Observer) OnCollation(c *collate.Collation) {
	serr := o.sparks.sparkCollation(c)
	lerr := o.logCollation(c)
	o.logError(iohelp.FirstError(serr, lerr))
}

func (o *Observer) logCollation(c *collate.Collation) error {
	sc := collate.Sourced{
		Run: run.Run{
			MachineID: o.mid,
			Iter:      o.nruns,
			Start:     o.lastTime,
		},
		Collation: c,
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
