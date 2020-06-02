// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package recipe_test

import (
	"fmt"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
	"github.com/MattWindsor91/act-tester/internal/model/recipe"
)

// ExampleOp_String is a runnable example for String.
func ExampleOp_String() {
	fmt.Println(recipe.CompileExe)
	fmt.Println(recipe.Op(42))

	// Output:
	// CompileExe
	// Op(42)
}

// TestOp_MarshalJSON_roundTrip tests Op's marshalling and unmarshalling by round-trip.
func TestOp_MarshalJSON_roundTrip(t *testing.T) {
	t.Parallel()
	for i := recipe.Nop; i <= recipe.Last; i++ {
		i := i
		t.Run(i.String(), func(t *testing.T) {
			testhelp.TestJSONRoundTrip(t, i, "round-trip Op")
		})
	}
}
