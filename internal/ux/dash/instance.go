// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	copy2 "github.com/MattWindsor91/c4t/internal/copier"
	"github.com/MattWindsor91/c4t/internal/director"
	"github.com/MattWindsor91/c4t/internal/observing"
	"github.com/MattWindsor91/c4t/internal/stage/analyser/saver"
	"github.com/MattWindsor91/c4t/internal/stage/mach/observer"
	"github.com/MattWindsor91/c4t/internal/stage/perturber"
	"github.com/MattWindsor91/c4t/internal/subject/status"

	"github.com/MattWindsor91/c4t/internal/plan/analysis"

	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"

	"github.com/MattWindsor91/c4t/internal/subject/corpus/builder"
)

const (
	headerCycle           = "Cycle"
	headerCurrentActivity = "Current Activity"
	headerStats           = "Statistics"
	headerSparks          = "Sparklines"
)

// Instance represents a single machine instance inside a dash.
type Instance struct {
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

// NewInstance constructs an Instance, initialising its various widgets.
// It accepts the id of the parent container (from which the IDs of various sub-containers can be derived), as well as
// the parent dash d (used to access the results log and parent container).
func NewInstance(id string, d *Dash) (*Instance, error) {
	o := Instance{
		id:   id,
		rlog: d.resultLog,
	}

	if err := o.setup(d); err != nil {
		return nil, err
	}

	return &o, nil
}

func (o *Instance) setup(d *Dash) error {
	var err error
	if o.tally, err = newTally(); err != nil {
		return err
	}
	if o.run, err = newRunCounter(); err != nil {
		return err
	}
	if o.sparks, err = newSparkset(); err != nil {
		return err
	}
	if o.action, err = NewCorpusObserver(); err != nil {
		return err
	}
	if o.compilers, err = newCompilerObserver(d.container, o.compilersContainerID()); err != nil {
		return err
	}
	return nil
}

const (
	percRun     = 25
	percStats   = 25
	percActions = 100 - percRun - percStats
)

// AddToGrid adds this observer to a grid builder gb with the container ID id..
func (o *Instance) AddToGrid(gb *grid.Builder, pc int) {
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

func (o *Instance) currentRunColumn() grid.Element {
	return grid.ColWidthPercWithOpts(percActions,
		[]container.Option{
			container.Border(linestyle.Light),
			container.BorderTitle(headerCurrentActivity),
		},
		o.action.gridRows()...,
	)
}

// OnIteration logs that a new iteration has begun.
func (o *Instance) OnCycle(r director.CycleMessage) {
	switch r.Kind {
	case director.CycleStart:
		o.onCycleStart(r.Cycle)
	}
}

func (o *Instance) onCycleStart(r director.Cycle) {
	o.nruns = r.Iter
	_ = o.run.onCycleStart(r)
	o.action.reset()
}

func (o *Instance) OnInstanceClose() {
	o.action.onInstanceClose()
}

// OnAnalysis observes an analysis by adding failure/timeout/flag rates to the sparklines.
func (o *Instance) OnAnalysis(a analysis.Analysis) {
	for i := status.Ok; i <= status.Last; i++ {
		o.sendStatusCount(i, len(a.ByStatus[i]))
	}
	if err := o.logAnalysis(a); err != nil {
		o.logError(err)
	}
}

// OnArchive currently ignores a save observation.
func (o *Instance) OnArchive(saver.ArchiveMessage) {
	// TODO(@MattWindsor91): do something with this?
}

func (o *Instance) sendStatusCount(i status.Status, n int) {
	if err := o.tally.tallyStatus(i, n); err != nil {
		o.logError(err)
	}
	if err := o.sparks.sparkStatus(i, n); err != nil {
		o.logError(err)
	}
}

func (o *Instance) logAnalysis(a analysis.Analysis) error {
	sc := director.CycleAnalysis{
		Cycle:    o.run.last,
		Analysis: a,
	}
	return o.rlog.Log(sc)
}

// OnBuild forwards a build observation.
func (o *Instance) OnBuild(m builder.Message) {
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

// OnCopy forwards a copy observation.
func (o *Instance) OnCopy(m copy2.Message) {
	switch m.Kind {
	case observing.BatchStart:
		o.action.OnCopyStart(m.Num)
	case observing.BatchStep:
		o.action.OnCopy(m.Dst, m.Src)
	case observing.BatchEnd:
		o.action.OnCopyFinish()
	}
}

// OnPerturb does nothing, at the moment.
func (o *Instance) OnPerturb(perturber.Message) {}

// OnMachineNodeAction does nothing, at the moment.
func (o *Instance) OnMachineNodeAction(observer.Message) {}

func (o *Instance) logError(err error) {
	// For want of better location.
	o.action.logError(err)
}
