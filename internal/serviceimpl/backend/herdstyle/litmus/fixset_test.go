// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus_test

import (
	"fmt"
	"os"

	"github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle/litmus"
)

// ExampleFixset_Args is a runnable example for Args.
func ExampleFixset_Args() {
	f := litmus.Fixset{InjectStdbool: true, UseAsCall: true}
	for _, s := range f.Args() {
		fmt.Println(s)
	}

	// Output:
	// -ascall
	// true
}

// ExampleFixset_Dump is a runnable example for Write.
func ExampleFixset_Dump() {
	f := litmus.Fixset{InjectStdbool: true, UseAsCall: true}
	_ = f.Dump(os.Stdout)

	// Output:
	// injecting stdbool
	// using -ascall
}

// ExampleFixset_NeedsPatch is a runnable example for NeedsPatch.
func ExampleFixset_NeedsPatch() {
	fmt.Println((&litmus.Fixset{}).NeedsPatch())
	fmt.Println((&litmus.Fixset{UseAsCall: true}).NeedsPatch())
	fmt.Println((&litmus.Fixset{InjectStdbool: true}).NeedsPatch())

	// Output:
	// false
	// false
	// true
}
