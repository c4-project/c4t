// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"

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

// OnStart sets up a DashObserver for a test phase with nreqs incoming requests.
func (d *Observer) OnStart(nreqs int) {
	d.nreqs = nreqs
	d.ndone = 0
	_ = d.g.Absolute(d.ndone, d.nreqs)
}

// OnAdd acknowledges the addition of a subject to a corpus being built.
func (d *Observer) OnAdd(sname string) {
	d.logAndStepGauge("ADD", sname, colorAdd)
}

// OnCompile acknowledges the addition of a compilation to a corpus being built.
func (d *Observer) OnCompile(sname string, compiler model.ID, success bool) {
	c := colorCompile
	desc := idQualSubjectDesc(sname, compiler)

	if !success {
		c = colorFailed
		desc += " [FAILED]"
	}

	d.logAndStepGauge("COMPILE", desc, c)
}

// OnHarness acknowledges the addition of a harness to a corpus being built.
func (d *Observer) OnHarness(sname string, arch model.ID) {
	d.logAndStepGauge("LIFT", idQualSubjectDesc(sname, arch), colorHarness)
}

// OnRun acknowledges the addition of a run to a corpus being built.
func (d *Observer) OnRun(sname string, compiler model.ID, s subject.Status) {
	desc := idQualSubjectDesc(sname, compiler)
	suff, c := runSuffixAndColour(s)
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
