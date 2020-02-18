package litmus

import (
	"fmt"
	"os"
)

// ExampleFixset_Dump is a runnable example for Dump.
func ExampleFixset_Dump() {
	f := Fixset{InjectStdbool: true, UseAsCall: true}
	_ = f.Dump(os.Stdout)

	// Output:
	// injecting stdbool
	// using -ascall
}

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
