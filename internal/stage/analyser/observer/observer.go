// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package observer

import "github.com/MattWindsor91/act-tester/internal/plan/analyser"

// Observer represents the observer interface for the analyser stage.
type Observer interface {
	// OnAnalysis lets the observer know that the current plan has been analysed and the results are in a.
	OnAnalysis(a analyser.Analysis)

	// OnArchive lets the observer know that an archive action has occurred.
	OnArchive(s ArchiveMessage)
}

//go:generate mockery -name=Observer

// OnAnalysis sends OnAnalysis to every instance observer in obs.
func OnAnalysis(a analyser.Analysis, obs ...Observer) {
	for _, o := range obs {
		o.OnAnalysis(a)
	}
}
