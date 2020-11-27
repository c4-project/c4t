// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package filekind_test

import (
	"fmt"

	"github.com/MattWindsor91/c4t/internal/model/filekind"
)

// ExampleGuessFromFile is a runnable example for GuessFromFile.
func ExampleGuessFromFile() {
	fmt.Println(filekind.GuessFromFile("foo.litmus.c") == filekind.CSrc)
	fmt.Println(filekind.GuessFromFile("stdio.h") == filekind.CHeader)
	fmt.Println(filekind.GuessFromFile("foo.litmus") == filekind.Litmus)
	fmt.Println(filekind.GuessFromFile("foo.litmus.c") == filekind.Litmus)

	// Output:
	// true
	// true
	// true
	// false
}

// ExampleKind_FilterFiles is a runnable example for FilterFiles.
func ExampleKind_FilterFiles() {
	for _, f := range filekind.C.FilterFiles([]string{
		"barbaz.trace",
		"foo.c",
		"foo.h",
		"baz.sh",
		"a.out",
		"bar.c.litmus",
		"bar.c",
	}) {
		fmt.Println(f)
	}

	// Output:
	// foo.c
	// foo.h
	// bar.c
}
