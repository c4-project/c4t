// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package csv_test

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/plan/analyser"
	acsv "github.com/MattWindsor91/act-tester/internal/stage/analyser/csv"
)

// NB: the below CSV is likely to change as the plan mock changes.
// At time of writing, the mock referred to compilers not in the plan, for instance.

// TODO(@MattWindsor91): add stages to the mock plan!

// ExampleStageWriter_OnAnalysis is a testable example for OnAnalysis.
func ExampleStageWriter_OnAnalysis() {
	az, _ := analyser.New(plan.Mock(), 1)
	an, _ := az.Analyse(context.Background())

	w := csv.NewWriter(os.Stdout)
	sw := (*acsv.StageWriter)(w)
	sw.OnAnalysis(*an)

	// Output:
	// Stage,CompletedAt,Duration
}
