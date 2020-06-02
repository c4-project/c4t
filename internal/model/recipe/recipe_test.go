// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package recipe_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"
)

// ExampleRecipe_Paths is a testable example for Paths.
func ExampleRecipe_Paths() {
	h := recipe.New("foo/bar", recipe.AddFiles("baz", "barbaz", "foobar"))
	for _, f := range h.Paths() {
		fmt.Println(f)
	}

	// Output:
	// foo/bar/baz
	// foo/bar/barbaz
	// foo/bar/foobar
}
