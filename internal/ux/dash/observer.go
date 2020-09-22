// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	copy2 "github.com/MattWindsor91/act-tester/internal/copier"
	"github.com/MattWindsor91/act-tester/internal/model/run"
	"github.com/MattWindsor91/act-tester/internal/observing"
	"github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer"
	"github.com/MattWindsor91/act-tester/internal/stage/perturber"
	"github.com/MattWindsor91/act-tester/internal/stage/planner"
	"github.com/MattWindsor91/act-tester/internal/subject/status"

	"github.com/MattWindsor91/act-tester/internal/plan/analysis"

	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"
)

const (
	headerCycle           = "Cycle"
	headerCurrentActivity = "Current Activity"
	headerStats           = "Statistics"
	headerSparks          = "Sparklines"
)

// Observer represents a single machine instance inside a dash.
type Observer struct {
	// id is the ID of the container for this machine observer, and also the base for deriving other subcontainer IDs.
	id string

	run   *runCounter
	rlog  *ResultLog
	tally *tally

	sparks *sparkset

	action    *actionObserver
	compilers *compilerObserver

	nruns uint64
}

// NewObserver constructs an Observer, initialising its various widgets.
// It accepts the id of the parent container (from which the IDs of various sub-containers can be derived), as well as
// the parent dash d (used to access the results log and parent container).
func NewObserver(id string, d *Dash) (*Observer, error) {
	var err error

	o := Observer{
		id:   id,
		rlog: d.resultLog,
	}

	if o.tally, err = newTally(); err != nil {
		return nil, err
	}

	if o.run, err = newRunCounter(); err != nil {
		return nil, err
	}

	if o.sparks, err = newSparkset(); err != nil {
		return nil, err
	}

	if o.action, err = NewCorpusObserver(); err != nil {
		return nil, err
	}

	if o.compilers, err = newCompilerObserver(d.container, o.compilersContainerID()); err != nil {
		return nil, err
	}

	return &o, nil
}

const (
	percRun     = 25
	percStats   = 25
	percActions = 100 - percRun - percStats
)

// AddToGrid adds this observer to a grid builder gb with the container ID id..
func (o *Observer) AddToGrid(gb *grid.Builder, pc int) {
	gb.Add(grid.RowHeightPercWithOpts(pc,
		[]container.Option{
			container.ID(o.id),
			container.Border(linestyle.Double),
		},
		grid.ColWidthPerc(percRun,
			grid.RowHeightPercWithOpts(
				40,
				[]container.Option{
					container.Border(linestyle.Light),
					container.BorderTitle(headerCycle),
				},
				o.run.grid()...,
			),
			grid.RowHeightPercWithOpts(
				60,
				[]container.Option{container.ID(o.compilersContainerID()), container.Border(linestyle.Light), container.BorderTitle(headerCompilers)},
				o.compilers.grid()...,
			),
		),
		grid.ColWidthPerc(percStats,
			grid.RowHeightPercWithOpts(
				40,
				[]container.Option{container.Border(linestyle.Light), container.BorderTitle(headerStats)},
				o.tally.grid()...,
			),
			grid.RowHeightPercWithOpts(60,
				[]container.Option{container.Border(linestyle.Light), container.BorderTitle(headerSparks)},
				o.sparks.grid()...),
		),
		o.currentRunColumn(),
	))
}

func (o *Observer) currentRunColumn() grid.Element {
	return grid.ColWidthPercWithOpts(percActions,
		[]container.Option{
			container.Border(linestyle.Light),
			container.BorderTitle(headerCurrentActivity),
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
	for i := status.Ok; i <= status.Last; i++ {
		o.sendStatusCount(i, len(a.ByStatus[i]))
	}
	if err := o.logAnalysis(a); err != nil {
		o.logError(err)
	}
}

// OnArchive currently ignores a save observation.
func (o *Observer) OnArchive(saver.ArchiveMessage) {
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
	sc := analysis.WithRun{
		Run:      o.run.last,
		Analysis: a,
	}
	return o.rlog.Log(sc)
}

// OnBuild forwards a build observation.
func (o *Observer) OnBuild(m builder.Message) {
	switch m.Kind {
	case observing.BatchStart:
		o.action.OnBuildStart(builder.Manifest{
			Name:  m.Name,
			NReqs: m.Num,
		})
	case observing.BatchStep:
		o.action.OnBuildRequest(*m.Request)
	case observing.BatchEnd:
		o.action.OnBuildFinish()
	}
}

// OnCopyStart forwards a copy start observation.
func (o *Observer) OnCopyStart(nfiles int) {
	o.action.OnCopyStart(nfiles)
}

// OnCopy forwards a copy observation.
func (o *Observer) OnCopy(m copy2.Message) {
	switch m.Kind {
	case observing.BatchStart:
		o.action.OnCopyStart(m.Num)
	case observing.BatchStep:
		o.action.OnCopy(m.Dst, m.Src)
	case observing.BatchEnd:
		o.action.OnCopyFinish()
	}
}

// OnPlan does nothing, at the moment.
func (o *Observer) OnPlan(planner.Message) {}

// OnPerturb does nothing, at the moment.
func (o *Observer) OnPerturb(perturber.Message) {}

// OnMachineNodeAction does nothing, at the moment.
func (o *Observer) OnMachineNodeAction(observer.Message) {}

func (o *Observer) logError(err error) {
	// For want of better location.
	o.action.logError(err)
}
