// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stage_test

import (
	"fmt"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
	"github.com/MattWindsor91/act-tester/internal/plan/stage"
)

// ExampleStage_String is a testable example for String.
func ExampleStage_String() {
	for i := stage.Unknown; i <= stage.Last+1; i++ {
		fmt.Println(i)
	}

	// Output:
	// Unknown
	// Plan
	// Fuzz
	// Lift
	// Invoke
	// Compile
	// Run
	// Analyse
	// Stage(8)
}

// ExampleStage_MarshalJSON is a runnable example for MarshalJSON.
func ExampleStage_MarshalJSON() {
	for i := stage.Unknown + 1; i <= stage.Last; i++ {
		bs, _ := i.MarshalJSON()
		fmt.Println(string(bs))
	}

	// Output:
	// "Plan"
	// "Fuzz"
	// "Lift"
	// "Invoke"
	// "Compile"
	// "Run"
	// "Analyse"
}

// TestStage_MarshalJSON_roundTrip tests Op's marshalling and unmarshalling by round-trip.
func TestStage_MarshalJSON_roundTrip(t *testing.T) {
	t.Parallel()
	for i := stage.Unknown; i <= stage.Last; i++ {
		i := i
		t.Run(i.String(), func(t *testing.T) {
			testhelp.TestJSONRoundTrip(t, i, "round-trip Stage")
		})
	}
}
