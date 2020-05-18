// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"
)

// ExampleAnalysis_String is a runnable example for String.
func ExampleAnalysis_String() {
	c := analysis.Analysis{
		ByStatus: map[status.Status]corpus.Corpus{
			status.Ok:             corpus.New("a", "b", "c", "ch"),
			status.Flagged:        corpus.New("barbaz"),
			status.CompileFail:    corpus.New("foo", "bar", "baz"),
			status.CompileTimeout: corpus.New(),
			status.RunFail:        corpus.New("foobaz", "barbaz"),
			status.RunTimeout:     corpus.New(),
		},
	}
	fmt.Println(&c)

	// Output:
	// 4 ok, 1 flagged, 3 compile/fail, 0 compile/timeout, 2 run/fail, 0 run/timeout
}

// ExampleAnalysis_HasFlagged is a runnable example for HasFailures.
func ExampleAnalysis_HasFlagged() {
	var empty analysis.Analysis
	fmt.Println("empty:", empty.HasFlagged())

	flagged := analysis.Analysis{
		ByStatus: map[status.Status]corpus.Corpus{
			status.Flagged: corpus.New("foo", "bar", "baz"),
		},
		Flags: analysis.FlagFlagged,
	}
	fmt.Println("flagged:", flagged.HasFlagged())

	// Output:
	// empty: false
	// flagged: true
}

// ExampleAnalysis_HasFailures is a runnable example for HasFailures.
func ExampleAnalysis_HasFailures() {
	var empty analysis.Analysis
	fmt.Println("empty:", empty.HasFailures())

	cfails := analysis.Analysis{
		Flags: analysis.FlagCompileFail,
	}
	fmt.Println("compiler failures:", cfails.HasFailures())

	rfails := analysis.Analysis{
		Flags: analysis.FlagRunFail,
	}
	fmt.Println("run failures:", rfails.HasFailures())

	// Output:
	// empty: false
	// compiler failures: true
	// run failures: true
}
