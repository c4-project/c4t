// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package analyser

import (
	"github.com/c4-project/c4t/internal/plan/analysis"
)

// Observer represents the observer interface for the analyser stage.
type Observer interface {
	// OnAnalysis lets the observer know that the current plan has been analysed and the results are in a.
	OnAnalysis(a analysis.Analysis)
}

//go:generate mockery --name=Observer

// OnAnalysis sends OnAnalysis to every instance observer in obs.
func OnAnalysis(a analysis.Analysis, obs ...Observer) {
	for _, o := range obs {
		o.OnAnalysis(a)
	}
}
