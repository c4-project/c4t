// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/collate"
)

// ExampleCollation_String is a runnable example for String.
func ExampleCollation_String() {
	c := collate.Collation{
		ByStatus: map[subject.Status]corpus.Corpus{
			subject.StatusOk:             corpus.New("a", "b", "c", "ch"),
			subject.StatusFlagged:        corpus.New("barbaz"),
			subject.StatusCompileFail:    corpus.New("foo", "bar", "baz"),
			subject.StatusCompileTimeout: corpus.New(),
			subject.StatusRunFail:        corpus.New("foobaz", "barbaz"),
			subject.StatusRunTimeout:     corpus.New(),
		},
	}
	fmt.Println(&c)

	// Output:
	// 4 ok, 1 flagged, 3 compile/fail, 0 compile/timeout, 2 run/fail, 0 run/timeout
}
