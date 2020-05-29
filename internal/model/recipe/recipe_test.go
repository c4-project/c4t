// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package recipe

import "fmt"

// ExampleRecipe_Paths is a testable example for Paths.
func ExampleRecipe_Paths() {
	h := Recipe{Dir: "foo/bar", Files: []string{"baz", "barbaz", "foobar"}}
	for _, f := range h.Paths() {
		fmt.Println(f)
	}

	// Output:
	// foo/bar/baz
	// foo/bar/barbaz
	// foo/bar/foobar
}
