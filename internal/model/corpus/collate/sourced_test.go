// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate_test

import (
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/collate"
	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// ExampleSourced_String is a runnable example for String.
func ExampleSourced_String() {
	sc := collate.Sourced{
		MachineID: id.FromString("foo.bar.baz"),
		Iter:      42,
		Start:     time.Date(1997, time.May, 1, 10, 0, 0, 0, time.FixedZone("BST", 60*60)),
		Collation: nil,
	}

	// Without collation:
	fmt.Println(&sc)

	// With collation:
	sc.Collation = &collate.Collation{
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
	fmt.Println(&sc)

	// Output:
	// [foo.bar.baz #42 (May  1 10:00:00)] (nil)
	// [foo.bar.baz #42 (May  1 10:00:00)] 4 ok, 1 flagged, 3 compile/fail, 0 compile/timeout, 2 run/fail, 0 run/timeout
}
