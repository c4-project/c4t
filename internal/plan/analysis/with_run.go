// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"github.com/MattWindsor91/act-tester/internal/model/run"
)

// WithRun contains a corpus collation and its parent run.
type WithRun struct {
	// Run contains information about the run that produced this collation.
	Run run.Run

	// Analysis is the collation proper.
	Analysis
}

// String formats a log header for this sourced analysis.
func (s *WithRun) String() string {
	return s.Run.String() + " " + s.Analysis.String()
}
