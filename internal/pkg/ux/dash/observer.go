// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mum4k/termdash/widgets/sparkline"

	"github.com/mum4k/termdash/widgets/segmentdisplay"

	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/gauge"
	"github.com/mum4k/termdash/widgets/text"
)

const (
	colorAdd     = cell.ColorBlue
	colorCompile = cell.ColorMagenta
	colorFailed  = cell.ColorRed
	colorFlagged = cell.ColorYellow
	colorHarness = cell.ColorCyan
	colorRun     = cell.ColorGreen
	colorTimeout = colorFailed
)

// Observer is a BuilderObserver that attaches into a Dash.
type Observer struct {
	runCount     *segmentdisplay.SegmentDisplay
	runStart     *text.Text
	runTimeSpark *sparkline.SparkLine
	buildLog     *text.Text
	buildGauge   *gauge.Gauge
	lastTime     time.Time
	nreqs, ndone int
}

// NewObserver constructs an Observer, initialising its various widgets.
func NewObserver() (*Observer, error) {
	var (
		d   Observer
		err error
	)

	if d.runCount, err = segmentdisplay.New(); err != nil {
		return nil, err
	}

	if d.runStart, err = text.New(text.DisableScrolling()); err != nil {
		return nil, err
	}

	if d.runTimeSpark, err = sparkline.New(); err != nil {
		return nil, err
	}

	if d.buildGauge, err = gauge.New(); err != nil {
		return nil, err
	}

	if d.buildLog, err = text.New(text.RollContent()); err != nil {
		return nil, err
	}

	return &d, nil
}

// AddToGrid adds this observer to a grid builder gb.
func (o *Observer) AddToGrid(gb *grid.Builder, midstr string, pc int) {
	gb.Add(grid.RowHeightPercWithOpts(pc,
		[]container.Option{container.Border(linestyle.Double), container.BorderTitle(midstr)},
		grid.ColWidthPerc(30,
			grid.RowHeightFixed(1, grid.Widget(o.runStart)),
			grid.RowHeightFixed(1, grid.Widget(o.runTimeSpark)),
			grid.RowHeightFixed(1, grid.Widget(o.runCount, container.Border(linestyle.Round), container.BorderTitle("Run#"))),
		),
		o.currentRunColumn(),
	))
}

func (o *Observer) currentRunColumn() grid.Element {
	return grid.ColWidthPercWithOpts(70,
		[]container.Option{
			container.Border(linestyle.Light),
			container.BorderTitle("Current Run"),
		},
		grid.RowHeightFixed(1, grid.Widget(o.buildGauge)),
		grid.RowHeightFixed(1, grid.Widget(o.buildLog)),
	)
}

// OnIteration logs that a new iteration has begun.
func (o *Observer) OnIteration(iter uint64, t time.Time) {
	ch := segmentdisplay.NewChunk(strconv.FormatUint(iter, 10))
	_ = o.runCount.Write(
		[]*segmentdisplay.TextChunk{ch},
	)

	_ = o.runStart.Write(t.Format(time.Stamp), text.WriteReplace())

	o.addDurationToSparkline(t)
}

func (o *Observer) addDurationToSparkline(t time.Time) {
	if !o.lastTime.IsZero() {
		dur := t.Sub(o.lastTime)
		_ = o.runTimeSpark.Add([]int{int(dur.Seconds())})
	}
	o.lastTime = t
}

// OnStart sets up an observer for a test phase with manifest m.
func (o *Observer) OnStart(m builder.Manifest) {
	_ = o.buildLog.Write(fmt.Sprintf("-- %s --\n", m.Name))

	o.nreqs = m.NReqs
	o.ndone = 0
	_ = o.buildGauge.Absolute(o.ndone, o.nreqs)
}

// OnRequest acknowledges a corpus-builder request.
func (o *Observer) OnRequest(r builder.Request) {
	switch {
	case r.Add != nil:
		o.onAdd(r.Name)
	case r.Compile != nil:
		o.onCompile(r.Name, r.Compile)
	case r.Harness != nil:
		o.onHarness(r.Name, r.Harness)
	case r.Run != nil:
		o.onRun(r.Name, r.Run)
	}
}

// onAdd acknowledges the addition of a subject to a corpus being built.
func (o *Observer) onAdd(sname string) {
	o.logAndStepGauge("ADD", sname, colorAdd)
}

// onCompile acknowledges the addition of a compilation to a corpus being built.
func (o *Observer) onCompile(sname string, b *builder.Compile) {
	c := colorCompile
	desc := idQualSubjectDesc(sname, b.CompilerID)

	if !b.Result.Success {
		c = colorFailed
		desc += " [FAILED]"
	}

	o.logAndStepGauge("COMPILE", desc, c)
}

// onHarness acknowledges the addition of a harness to a corpus being built.
func (o *Observer) onHarness(sname string, b *builder.Harness) {
	o.logAndStepGauge("LIFT", idQualSubjectDesc(sname, b.Arch), colorHarness)
}

// onRun acknowledges the addition of a run to a corpus being built.
func (o *Observer) onRun(sname string, b *builder.Run) {
	desc := idQualSubjectDesc(sname, b.CompilerID)
	suff, c := runSuffixAndColour(b.Result.Status)
	o.logAndStepGauge("RUN", desc+suff, c)
}

func runSuffixAndColour(s subject.Status) (string, cell.Color) {
	switch s {
	case subject.StatusFlagged:
		return " [FLAGGED]", colorFlagged
	case subject.StatusTimeout:
		return " [TIMEOUT]", colorTimeout
	case subject.StatusCompileFail:
		return " [FAILED]", colorFailed
	default:
		return "", colorRun
	}
}

// OnFinish acknowledges the end of a run.
func (o *Observer) OnFinish() {
	_ = o.buildLog.Write("-- DONE --\n")
}

func idQualSubjectDesc(sname string, id id.ID) string {
	return fmt.Sprintf("%s (@%s)", sname, id)
}

// logAndStepGauge logs a request with name rq and summary desc, then repopulates the gauge.
// It uses c as the colour for both.
func (o *Observer) logAndStepGauge(rq, desc string, c cell.Color) {
	o.log(rq, desc, c)
	o.stepGauge(c)
}

// log logs an observed builder request with name rq and summary desc to the per-machine log.
// It colours the log with c.
func (o *Observer) log(rq, desc string, c cell.Color) {
	_ = o.buildLog.Write(rq, text.WriteCellOpts(cell.FgColor(c)))
	_ = o.buildLog.Write(" " + desc + "\n")
}

// stepGauge increments the gauge and sets its colour to c.
func (o *Observer) stepGauge(c cell.Color) {
	o.ndone++
	_ = o.buildGauge.Absolute(o.ndone, o.nreqs, gauge.Color(c))
}
