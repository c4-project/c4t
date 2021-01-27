// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"fmt"
	"os"
)

// ExampleFixset_Args is a runnable example for Args.
func ExampleFixset_Args() {
	f := Fixset{InjectStdbool: true, UseAsCall: true}
	for _, s := range f.Args() {
		fmt.Println(s)
	}

	// Output:
	// -ascall
	// true
}

// ExampleFixset_Dump is a runnable example for Write.
func ExampleFixset_Dump() {
	f := Fixset{InjectStdbool: true, UseAsCall: true}
	_ = f.Dump(os.Stdout)

	// Output:
	// injecting stdbool
	// using -ascall
}

// ExampleFixset_NeedsPatch is a runnable example for NeedsPatch.
func ExampleFixset_NeedsPatch() {
	fmt.Println((&Fixset{}).NeedsPatch())
	fmt.Println((&Fixset{UseAsCall: true}).NeedsPatch())
	fmt.Println((&Fixset{InjectStdbool: true}).NeedsPatch())

	// Output:
	// false
	// false
	// true
}
