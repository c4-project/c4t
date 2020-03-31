// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/collate"
)

// ExampleCollation_HasFlagged is a runnable example for HasFailures.
func ExampleCollation_HasFlagged() {
	var empty collate.Collation
	fmt.Println("empty:", empty.HasFlagged())

	flagged := collate.Collation{
		Flagged: corpus.New("foo", "bar", "baz"),
	}
	fmt.Println("flagged:", flagged.HasFlagged())

	// Output:
	// empty: false
	// flagged: true
}

// ExampleCollation_HasFailures is a runnable example for HasFailures.
func ExampleCollation_HasFailures() {
	var empty collate.Collation
	fmt.Println("empty:", empty.HasFailures())

	cfails := collate.Collation{
		Compile: collate.FailCollation{
			Failures: corpus.New("foo", "bar", "baz"),
		},
	}
	fmt.Println("compiler failures:", cfails.HasFailures())

	rfails := collate.Collation{
		Run: collate.FailCollation{
			Failures: corpus.New("foo", "bar", "baz"),
		},
	}
	fmt.Println("run failures:", rfails.HasFailures())

	// Output:
	// empty: false
	// compiler failures: true
	// run failures: true
}

// ExampleFailCollation_IsEmpty is a runnable example for IsEmpty.
func ExampleFailCollation_IsEmpty() {
	var empty collate.FailCollation
	fmt.Println("empty:", empty.IsEmpty())

	notEmpty := collate.FailCollation{
		Failures: corpus.New("foo", "bar", "baz"),
	}
	fmt.Println("not empty:", notEmpty.IsEmpty())

	// Output:
	// empty: true
	// not empty: false
}
