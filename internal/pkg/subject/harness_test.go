// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject

import "fmt"

// ExampleHarness_Paths is a testable example for Paths.
func ExampleHarness_Paths() {
	h := Harness{Dir: "foo/bar", Files: []string{"baz", "barbaz", "foobar"}}
	for _, f := range h.Paths() {
		fmt.Println(f)
	}

	// Output:
	// foo/bar/baz
	// foo/bar/barbaz
	// foo/bar/foobar
}

// ExampleHarness_CPaths is a testable example for CPaths.
func ExampleHarness_CPaths() {
	h := Harness{Dir: "foo/bar", Files: []string{"README.md", "baz.c", "barbaz.c", "foobar.h"}}
	for _, f := range h.CPaths() {
		fmt.Println(f)
	}

	// Output:
	// foo/bar/baz.c
	// foo/bar/barbaz.c
}
