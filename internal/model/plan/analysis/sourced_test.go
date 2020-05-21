// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis_test

import (
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/run"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"
)

// ExampleSourced_String is a runnable example for String.
func ExampleSourced_String() {
	sc := analysis.Sourced{
		Run: run.Run{
			MachineID: id.FromString("foo.bar.baz"),
			Iter:      42,
			Start:     time.Date(1997, time.May, 1, 10, 0, 0, 0, time.FixedZone("BST", 60*60)),
		},
		Analysis: analysis.Analysis{
			ByStatus: map[status.Status]corpus.Corpus{
				status.Ok:             corpus.New("a", "b", "c", "ch"),
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
	// [foo.bar.baz #42 (May  1 10:00:00)] 4 Ok, 1 Flagged, 3 CompileFail, 0 CompileTimeout, 2 RunFail, 0 RunTimeout
}
