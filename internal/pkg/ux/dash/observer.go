// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
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
	mid          model.ID
	last         *text.Text
	g            *gauge.Gauge
	nreqs, ndone int
}

// OnStart sets up a DashObserver for a test phase with manifest m
func (d *Observer) OnStart(m builder.Manifest) {
	// TODO(@MattWindsor91): use name
	_ = d.last.Write(fmt.Sprintf("-- %s --\n", m.Name))

	d.nreqs = m.NReqs
	d.ndone = 0
	_ = d.g.Absolute(d.ndone, d.nreqs)
}

// OnRequest acknowledges a corpus-builder request.
func (d *Observer) OnRequest(r builder.Request) {
	switch {
	case r.Add != nil:
		d.onAdd(r.Name)
	case r.Compile != nil:
		d.onCompile(r.Name, r.Compile)
	case r.Harness != nil:
		d.onHarness(r.Name, r.Harness)
	case r.Run != nil:
		d.onRun(r.Name, r.Run)
	}
}

// onAdd acknowledges the addition of a subject to a corpus being built.
func (d *Observer) onAdd(sname string) {
	d.logAndStepGauge("ADD", sname, colorAdd)
}

// onCompile acknowledges the addition of a compilation to a corpus being built.
func (d *Observer) onCompile(sname string, b *builder.Compile) {
	c := colorCompile
	desc := idQualSubjectDesc(sname, b.CompilerID)

	if !b.Result.Success {
		c = colorFailed
		desc += " [FAILED]"
	}

	d.logAndStepGauge("COMPILE", desc, c)
}

// onHarness acknowledges the addition of a harness to a corpus being built.
func (d *Observer) onHarness(sname string, b *builder.Harness) {
	d.logAndStepGauge("LIFT", idQualSubjectDesc(sname, b.Arch), colorHarness)
}

// onRun acknowledges the addition of a run to a corpus being built.
func (d *Observer) onRun(sname string, b *builder.Run) {
	desc := idQualSubjectDesc(sname, b.CompilerID)
	suff, c := runSuffixAndColour(b.Result.Status)
	d.logAndStepGauge("RUN", desc+suff, c)
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
func (d *Observer) OnFinish() {
	_ = d.last.Write("-- DONE --\n")
}

func idQualSubjectDesc(sname string, id model.ID) string {
	return fmt.Sprintf("%s (@%s)", sname, id)
}

// logAndStepGauge logs a request with name rq and summary desc, then repopulates the gauge.
// It uses c as the colour for both.
func (d *Observer) logAndStepGauge(rq, desc string, c cell.Color) {
	d.log(rq, desc, c)
	d.stepGauge(c)
}

// log logs an observed builder request with name rq and summary desc to the per-machine log.
// It colours the log with c.
func (d *Observer) log(rq, desc string, c cell.Color) {
	_ = d.last.Write(rq, text.WriteCellOpts(cell.FgColor(c)))
	_ = d.last.Write(" " + desc + "\n")
}

// stepGauge increments the gauge and sets its colour to c.
func (d *Observer) stepGauge(c cell.Color) {
	d.ndone++
	_ = d.g.Absolute(d.ndone, d.nreqs, gauge.Color(c))
}
