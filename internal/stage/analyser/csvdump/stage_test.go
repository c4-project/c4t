// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package csvdump_test

import (
	"context"
	"os"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/csvdump"

	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/plan/analysis"
)

// NB: the below CSV is likely to change as the plan mock changes.
// At time of writing, the mock referred to compilers not in the plan, for instance.

// TODO(@MattWindsor91): add stages to the mock plan!

// ExampleStageWriter_OnAnalysis is a testable example for OnAnalysis.
func ExampleStageWriter_OnAnalysis() {
	an, _ := analysis.Analyse(context.Background(), plan.Mock())

	sw := csvdump.NewStageWriter(os.Stdout)
	sw.OnAnalysis(*an)

	// Output:
	// Stage,CompletedAt,Duration
}
