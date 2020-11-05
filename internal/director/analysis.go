// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import "github.com/MattWindsor91/act-tester/internal/plan/analysis"

// CycleAnalysis contains an analysis as well as the cycle that produced it.
type CycleAnalysis struct {
	// Cycle contains information about the run that produced this collation.
	Cycle Cycle

	// Analysis is the collation proper.
	analysis.Analysis
}

// String formats a log header for this sourced analysis.
func (s *CycleAnalysis) String() string {
	return s.Cycle.String() + " " + s.Analysis.String()
}