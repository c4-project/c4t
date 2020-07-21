// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stage_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/plan/stage"
)

// ExampleRecord_String is a testable example for String.
func ExampleRecord_String() {
	for i := stage.Unknown + 1; i <= stage.Analyse+1; i++ {
		fmt.Println(i)
	}

	// Output:
	// Plan
	// Fuzz
	// Lift
	// Invoke
	// Analyse
	// Stage(6)
}
