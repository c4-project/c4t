// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"github.com/MattWindsor91/c4t/internal/director"
	"github.com/MattWindsor91/c4t/internal/model/service/compiler/optlevel"
	"github.com/MattWindsor91/c4t/internal/subject/status"
	"github.com/mum4k/termdash/cell"
)

const (
	colourCopy = cell.ColorWhite

	colourOptNone   = cell.ColorBlue
	colourOptNormal = cell.ColorMagenta
	colourOptBreak  = cell.ColorRed

	colourAdd  = cell.ColorBlue
	colourLift = cell.ColorCyan
	colourRun  = cell.ColorGreen

	colourUnknown        = cell.ColorWhite
	colourOk             = cell.ColorGreen
	colourFiltered       = cell.ColorWhite // colourUnknown is unlikely to appear in practice, so duplication is ok
	colourFlagged        = cell.ColorYellow
	colourCompileFail    = cell.ColorRed
	colourCompileTimeout = cell.ColorBlue
	colourRunFail        = cell.ColorMagenta
	colourRunTimeout     = cell.ColorCyan
)

// statusColours maps each status flag to its colour.
// This will need to be kept in sync with the status enum.
var statusColours = [status.Last + 1]cell.Color{
	colourUnknown,
	colourOk,
	colourFiltered,
	colourFlagged,
	colourCompileFail,
	colourCompileTimeout,
	colourRunFail,
	colourRunTimeout,
}

// optColour divines a colour to signify the optimisation level described by o.
func optColour(o *optlevel.Named) cell.Color {
	switch {
	case o == nil || !o.Optimises:
		return colourOptNone
	case o.BreaksStandards:
		return colourOptBreak
	default:
		return colourOptNormal
	}
}

// summaryColor retrieves a colour to use for the log header of sc, according to a 'traffic lights' system.
func summaryColor(sc director.CycleAnalysis) cell.Color {
	switch {
	case sc.Analysis.HasFailures():
		return colourCompileFail
	case sc.Analysis.HasFlagged():
		return colourFlagged
	default:
		return colourRun
	}
}
