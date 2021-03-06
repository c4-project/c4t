// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package quantity_test

import (
	"fmt"

	"github.com/c4-project/c4t/internal/quantity"
)

// ExampleFuzzSet_Override is a runnable example for FuzzSet.Override.
func ExampleFuzzSet_Override() {
	q1 := quantity.FuzzSet{
		CorpusSize:    27,
		SubjectCycles: 53,
	}
	q2 := quantity.FuzzSet{
		SubjectCycles: 42,
	}
	q1.Override(q2)

	fmt.Println("corpus size:   ", q1.CorpusSize)
	fmt.Println("subject cycles:", q1.SubjectCycles)

	// Output:
	// corpus size:    27
	// subject cycles: 42
}
