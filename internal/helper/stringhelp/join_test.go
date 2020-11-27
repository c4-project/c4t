// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stringhelp_test

import (
	"fmt"

	"github.com/MattWindsor91/c4t/internal/helper/stringhelp"
)

// ExampleJoinNonEmpty is a testable example for JoinNonEmpty.
func ExampleJoinNonEmpty() {
	fmt.Println(stringhelp.JoinNonEmpty("/"))
	fmt.Println(stringhelp.JoinNonEmpty("/", ""))
	fmt.Println(stringhelp.JoinNonEmpty("/", "example1", ""))
	fmt.Println(stringhelp.JoinNonEmpty("/", "", "example2"))
	fmt.Println(stringhelp.JoinNonEmpty("/", "example1", "example2"))
	fmt.Println(stringhelp.JoinNonEmpty("/", "the", "", "quick brown", "", "fox"))

	// Output:
	//
	//
	// example1
	// example2
	// example1/example2
	// the/quick brown/fox
}
