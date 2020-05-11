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

// ExampleCollation_HasFlagged is a runnable example for HasFailures.
func ExampleCollation_HasFlagged() {
	var empty collate.Collation
	fmt.Println("empty:", empty.HasFlagged())

	flagged := collate.Collation{
		ByStatus: map[subject.Status]corpus.Corpus{
			subject.StatusFlagged: corpus.New("foo", "bar", "baz"),
		},
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
		ByStatus: map[subject.Status]corpus.Corpus{
			subject.StatusCompileFail: corpus.New("foo", "bar", "baz"),
		},
	}
	fmt.Println("compiler failures:", cfails.HasFailures())

	rfails := collate.Collation{
		ByStatus: map[subject.Status]corpus.Corpus{
			subject.StatusRunFail: corpus.New("foo", "bar", "baz"),
		},
	}
	fmt.Println("run failures:", rfails.HasFailures())

	// Output:
	// empty: false
	// compiler failures: true
	// run failures: true
}
