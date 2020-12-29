// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"time"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/subject/obs"
)

// Result is the type of results from a single test run.
type Result struct {
	// Time is the time at which the run commenced.
	Time time.Time `json:"time,omitempty"`

	// Normalise is the set of results categorised by subject.
	// Each key is the subject's name.
	Subjects map[string]SubjectResult `json:"subject,omitempty"`
}

// SubjectResult contains results from a single subject.
type SubjectResult struct {
	// Compilers is the set of per-compiler results that were reported for this subject.
	// Each key is a stringified form of a compiler CompilerID.
	Compilers map[string]CompilerResult `json:"compiler,omitempty"`
}

// CompilerResult contains results from a subject/compiler pairing.
type CompilerResult struct {
	// Status is the status of this run.
	Status status.Status `json:"status"`

	// Obs is this subject's processed observation, if any.
	Obs *obs.Obs `json:"obs,omitempty"`
}
