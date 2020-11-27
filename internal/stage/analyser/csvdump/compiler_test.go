// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package csvdump_test

import (
	"context"
	"os"

	"github.com/MattWindsor91/c4t/internal/plan"
	"github.com/MattWindsor91/c4t/internal/plan/analysis"
	"github.com/MattWindsor91/c4t/internal/stage/analyser/csvdump"
)

// NB: the below CSV is likely to change as the plan mock changes.
// At time of writing, the mock referred to compilers not in the plan, for instance.

// ExampleCompilerWriter_OnAnalysis is a testable example for CompilerWriter.OnAnalysis.
func ExampleCompilerWriter_OnAnalysis() {
	an, _ := analysis.Analyse(context.Background(), plan.Mock())

	// nb: aside from the header, the actual order of compilers is not deterministic
	cw := csvdump.NewCompilerWriter(os.Stdout)
	cw.OnAnalysis(*an)

	// Unordered output:
	// CompilerID,StyleID,ArchID,Opt,MOpt,MinCompile,AvgCompile,MaxCompile,MinRun,AvgRun,MaxRun,Ok,Filtered,Flagged,CompileFail,CompileTimeout,RunFail,RunTimeout
	// gcc,gcc,ppc.64le.power9,,,200,200,200,0,0,0,0,0,1,1,0,0,0
	// clang,gcc,x86,,,200,200,200,0,0,0,1,0,0,0,0,0,0
}
