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

// TODO(@MattWindsor91): add things to the plan mock so that the below analysis gets populated.

// ExampleCompilerWriter_OnAnalysis is a testable example for OnAnalysis.
func ExampleCompilerWriter_OnAnalysis() {
	az, _ := analyser.New(plan.Mock(), 1)
	an, _ := az.Analyse(context.Background())

	w := csv.NewWriter(os.Stdout)
	cw := (*acsv.CompilerWriter)(w)
	cw.OnAnalysis(*an)

	// Output:
	// CompilerID,StyleID,ArchID,Opt,MOpt,MinCompile,AvgCompile,MaxCompile,MinRun,AvgRun,MaxRun,Ok,Flagged,CompileFail,CompileTimeout,RunFail,RunTimeout
}
