// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package obs_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/obs"
)

// ExampleFlagOfStrings is a testable example for ObsFlagOfStrings.
func ExampleFlagOfStrings() {
	f, _ := obs.FlagOfStrings("unsat", "undef")
	fmt.Println(f.Has(obs.Sat))
	fmt.Println(f.Has(obs.Unsat))
	fmt.Println(f.Has(obs.Undef))

	// Output:
	// false
	// true
	// true
}

// ExampleFlag_Strings is a testable example for Strings.
func ExampleFlag_Strings() {
	for _, s := range (obs.Sat | obs.Undef).Strings() {
		fmt.Println(s)
	}

	// Output:
	// sat
	// undef
}
