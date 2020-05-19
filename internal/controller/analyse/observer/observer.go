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

	// OnSave lets the observer know that a save action has occurred.
	OnSave(s Saving)

	// OnSaveFileMissing lets the observer know that a save action couldn't find a file.
	OnSaveFileMissing(s Saving, missing string)
}

// OnAnalysis sends OnAnalysis to every instance observer in obs.
func OnAnalysis(a analysis.Analysis, obs ...Observer) {
	for _, o := range obs {
		o.OnAnalysis(a)
	}
}

// OnSave sends OnSave to every instance observer in obs.
func OnSave(s Saving, obs ...Observer) {
	for _, o := range obs {
		o.OnSave(s)
	}
}

// OnSaveFileMissing sends OnSaveFileMissing to every instance observer in obs.
func OnSaveFileMissing(s Saving, missing string, obs ...Observer) {
	for _, o := range obs {
		o.OnSaveFileMissing(s, missing)
	}
}

// Saving represents a pair of saved subject and destination, used in OnSave messages.
type Saving struct {
	// SubjectName is the name of the subject that was saved
	SubjectName string
	// Dest is the destination of the saving action.
	Dest string
}
