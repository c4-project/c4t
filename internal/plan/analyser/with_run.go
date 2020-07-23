// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analyser

import (
	"github.com/MattWindsor91/act-tester/internal/model/run"
)

// AnalysisWithRun contains a corpus collation and its parent run.
type AnalysisWithRun struct {
	// Run contains information about the run that produced this collation.
	Run run.Run

	// Analysis is the collation proper.
	Analysis
}

// String formats a log header for this sourced collation.
func (s *AnalysisWithRun) String() string {
	return s.Run.String() + " " + s.Analysis.String()
}
