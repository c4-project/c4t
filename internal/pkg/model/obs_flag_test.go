// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package model_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// ExampleObsFlagOfStrings is a testable example for ObsFlagOfStrings.
func ExampleObsFlagOfStrings() {
	f, _ := model.ObsFlagOfStrings("unsat", "undef")
	fmt.Println(f.Has(model.ObsSat))
	fmt.Println(f.Has(model.ObsUnsat))
	fmt.Println(f.Has(model.ObsUndef))

	// Output:
	// false
	// true
	// true
}

// ExampleObsFlag_Strings is a testable example for Strings.
func ExampleObsFlag_Strings() {
	for _, s := range (model.ObsSat | model.ObsUndef).Strings() {
		fmt.Println(s)
	}

	// Output:
	// sat
	// undef
}
