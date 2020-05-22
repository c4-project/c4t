// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package observer

import "github.com/MattWindsor91/act-tester/internal/model/plan/analysis"

// Observer represents the observer interface for the analyse controller.
type Observer interface {
	// OnAnalysis lets the observer know that the current plan has been analysed and the results are in a.
	OnAnalysis(a analysis.Analysis)

	// OnArchive lets the observer know that an archive action has occurred.
	OnArchive(s ArchiveMessage)
}

//go:generate mockery -name=Observer

// OnAnalysis sends OnAnalysis to every instance observer in obs.
func OnAnalysis(a analysis.Analysis, obs ...Observer) {
	for _, o := range obs {
		o.OnAnalysis(a)
	}
}
