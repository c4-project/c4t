// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/helpers/iohelp"
)

// ExampleExtlessFile is a runnable example for ExtlessFile.
func ExampleExtlessFile() {
	fmt.Println(iohelp.ExtlessFile("foo.c"))
	fmt.Println(iohelp.ExtlessFile("/home/piers/test"))
	fmt.Println(iohelp.ExtlessFile("/home/piers/example.txt"))

	// Output:
	// foo
	// test
	// example
}
