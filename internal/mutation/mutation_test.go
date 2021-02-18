// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation_test

import (
	"fmt"

	"github.com/c4-project/c4t/internal/mutation"
)

// ExampleMutant_String is a runnable example for Mutant.String.
func ExampleMutant_String() {
	fmt.Println(mutation.NamedMutant(42, "", 0))
	fmt.Println(mutation.NamedMutant(56, "", 10))
	fmt.Println(mutation.NamedMutant(12, "FOO", 0))
	fmt.Println(mutation.NamedMutant(13, "BAR", 1))

	// Output:
	// 42
	// 56
	// FOO:12
	// BAR1:13
}
