// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/collate"
)

// ExampleCollation_ByStatus is a runnable example for ByStatus.
func ExampleCollation_ByStatus() {
	c := collate.Collation{
		Successes: corpus.New("a", "b", "c", "ch"),
		Flagged:   corpus.New("barbaz"),
		Compile: collate.FailCollation{
			Failures: corpus.New("foo", "bar", "baz"),
			Timeouts: corpus.New(),
		},
		Run: collate.FailCollation{
			Failures: corpus.New("foobaz", "barbaz"),
			Timeouts: corpus.New(),
		},
	}
	for k, v := range c.ByStatus() {
		fmt.Printf("%s:", k)
		for _, n := range v.Names() {
			fmt.Printf(" %s", n)
		}
		fmt.Println()
	}

	// Unordered output:
	// ok: a b c ch
	// flagged: barbaz
	// compile/fail: bar baz foo
	// compile/timeout:
	// run/fail: barbaz foobaz
	// run/timeout:
}

// ExampleCollation_String is a runnable example for String.
func ExampleCollation_String() {
	c := collate.Collation{
		Successes: corpus.New("a", "b", "c", "ch"),
		Flagged:   corpus.New("barbaz"),
		Compile: collate.FailCollation{
			Failures: corpus.New("foo", "bar", "baz"),
			Timeouts: corpus.New(),
		},
		Run: collate.FailCollation{
			Failures: corpus.New("foobaz", "barbaz"),
			Timeouts: corpus.New(),
		},
	}
	fmt.Println(&c)

	// Output:
	// 4 ok, 1 flagged, 3 compile/fail, 0 compile/timeout, 2 run/fail, 0 run/timeout
}
