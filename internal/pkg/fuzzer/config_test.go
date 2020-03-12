// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"
)

// ExampleQuantitySet_Override is a runnable example for Override.
func ExampleQuantitySet_Override() {
	q1 := fuzzer.QuantitySet{
		CorpusSize:    27,
		SubjectCycles: 53,
	}
	q2 := fuzzer.QuantitySet{
		SubjectCycles: 42,
	}
	q1.Override(q2)

	fmt.Println("corpus size:   ", q1.CorpusSize)
	fmt.Println("subject cycles:", q1.SubjectCycles)

	// Output:
	// corpus size:    27
	// subject cycles: 42
}
