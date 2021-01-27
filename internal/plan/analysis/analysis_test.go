// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package analysis_test

import (
	"fmt"

	"github.com/c4-project/c4t/internal/plan/analysis"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/subject/corpus"
)

// ExampleAnalysis_String is a runnable example for String.
func ExampleAnalysis_String() {
	c := analysis.Analysis{
		ByStatus: map[status.Status]corpus.Corpus{
			status.Ok:             corpus.New("a", "b", "c", "ch"),
			status.Filtered:       corpus.New("a", "i", "u", "e", "o"),
			status.Flagged:        corpus.New("barbaz"),
			status.CompileFail:    corpus.New("foo", "bar", "baz"),
			status.CompileTimeout: corpus.New(),
			status.RunFail:        corpus.New("foobaz", "barbaz"),
			status.RunTimeout:     corpus.New(),
		},
	}
	fmt.Println(&c)

	// Output:
	// 4 Ok, 5 Filtered, 1 Flagged, 3 CompileFail, 2 RunFail
}

// ExampleAnalysis_HasFlagged is a runnable example for Analysis.HasFlagged.
func ExampleAnalysis_HasFlagged() {
	var empty analysis.Analysis
	fmt.Println("empty:", empty.HasFlagged())

	flagged := analysis.Analysis{
		ByStatus: map[status.Status]corpus.Corpus{
			status.Flagged: corpus.New("foo", "bar", "baz"),
		},
		Flags: status.FlagFlagged,
	}
	fmt.Println("flagged:", flagged.HasFlagged())

	// Output:
	// empty: false
	// flagged: true
}

// ExampleAnalysis_HasFailures is a runnable example for Analysis.HasFailures.
func ExampleAnalysis_HasFailures() {
	var empty analysis.Analysis
	fmt.Println("empty:", empty.HasFailures())

	cfails := analysis.Analysis{Flags: status.FlagCompileFail}
	fmt.Println("compiler failures:", cfails.HasFailures())

	rfails := analysis.Analysis{Flags: status.FlagRunFail}
	fmt.Println("run failures:", rfails.HasFailures())

	ctos := analysis.Analysis{Flags: status.FlagCompileTimeout}
	fmt.Println("compiler timeouts:", ctos.HasFailures())

	rtos := analysis.Analysis{Flags: status.FlagRunTimeout}
	fmt.Println("run timeouts:", rtos.HasFailures())

	flags := analysis.Analysis{Flags: status.FlagFlagged}
	fmt.Println("flagged:", flags.HasFailures())

	filts := analysis.Analysis{Flags: status.FlagFiltered}
	fmt.Println("filtered:", filts.HasFailures())

	// Output:
	// empty: false
	// compiler failures: true
	// run failures: true
	// compiler timeouts: false
	// run timeouts: false
	// flagged: false
	// filtered: false
}

// ExampleAnalysis_HasBadOutcomes is a runnable example for Analysis.HasBadOutcomes.
func ExampleAnalysis_HasBadOutcomes() {
	var empty analysis.Analysis
	fmt.Println("empty:", empty.HasBadOutcomes())

	cfails := analysis.Analysis{Flags: status.FlagCompileFail}
	fmt.Println("compiler failures:", cfails.HasBadOutcomes())

	rfails := analysis.Analysis{Flags: status.FlagRunFail}
	fmt.Println("run failures:", rfails.HasBadOutcomes())

	ctos := analysis.Analysis{Flags: status.FlagCompileTimeout}
	fmt.Println("compiler timeouts:", ctos.HasBadOutcomes())

	rtos := analysis.Analysis{Flags: status.FlagRunTimeout}
	fmt.Println("run timeouts:", rtos.HasBadOutcomes())

	flags := analysis.Analysis{Flags: status.FlagFlagged}
	fmt.Println("flagged:", flags.HasBadOutcomes())

	filts := analysis.Analysis{Flags: status.FlagFiltered}
	fmt.Println("filtered:", filts.HasBadOutcomes())

	// Output:
	// empty: false
	// compiler failures: true
	// run failures: true
	// compiler timeouts: true
	// run timeouts: true
	// flagged: true
	// filtered: false
}
