// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package csvdump_test

import (
	"context"
	"os"

	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/plan/analyser"
	"github.com/MattWindsor91/act-tester/internal/stage/analyser/csvdump"
)

// NB: the below CSV is likely to change as the plan mock changes.
// At time of writing, the mock referred to compilers not in the plan, for instance.

// ExampleCompilerWriter_OnAnalysis is a testable example for OnAnalysis.
func ExampleCompilerWriter_OnAnalysis() {
	az, _ := analyser.New(plan.Mock(), 1)
	an, _ := az.Analyse(context.Background())

	// nb: aside from the header, the actual order of compilers is not deterministic
	cw := csvdump.NewCompilerWriter(os.Stdout)
	cw.OnAnalysis(*an)

	// Unordered output:
	// CompilerID,StyleID,ArchID,Opt,MOpt,MinCompile,AvgCompile,MaxCompile,MinRun,AvgRun,MaxRun,Ok,Flagged,CompileFail,CompileTimeout,RunFail,RunTimeout
	// gcc,gcc,ppc.64le.power9,,,200,200,200,0,0,0,0,1,1,0,0,0
	// clang,gcc,x86,,,200,200,200,0,0,0,1,0,0,0,0,0
}
