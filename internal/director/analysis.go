// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package director

import "github.com/c4-project/c4t/internal/plan/analysis"

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
