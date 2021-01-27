// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package director_test

import (
	"fmt"
	"time"

	"github.com/c4-project/c4t/internal/director"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/plan/analysis"
	"github.com/c4-project/c4t/internal/subject/corpus"
)

// ExampleCycleAnalysis_String is a runnable example for CycleAnalysis.String.
func ExampleCycleAnalysis_String() {
	sc := director.CycleAnalysis{
		Cycle: director.Cycle{
			Instance:  4,
			MachineID: id.FromString("foo.bar.baz"),
			Iter:      42,
			Start:     time.Date(1997, time.May, 1, 10, 0, 0, 0, time.FixedZone("BST", 60*60)),
		},
		Analysis: analysis.Analysis{
			ByStatus: map[status.Status]corpus.Corpus{
				status.Ok:             corpus.New("a", "b", "c", "ch"),
				status.Filtered:       corpus.New("a", "i", "u", "e", "o"),
				status.Flagged:        corpus.New("barbaz"),
				status.CompileFail:    corpus.New("foo", "bar", "baz"),
				status.CompileTimeout: corpus.New(),
				status.RunFail:        corpus.New("foobaz", "barbaz"),
				status.RunTimeout:     corpus.New(),
			},
		},
	}
	fmt.Println(&sc)

	// Output:
	// [4: foo.bar.baz #42 (May  1 10:00:00)] 4 Ok, 5 Filtered, 1 Flagged, 3 CompileFail, 2 RunFail
}
