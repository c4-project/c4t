// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import (
	"errors"
	"fmt"
)

// ExampleFirstError is a runnable example for FirstError.
func ExampleFirstError() {
	fmt.Println("FirstError() == nil:", FirstError() == nil)

	fmt.Println("FirstError(x, y) ==", FirstError(errors.New("x"), errors.New("y")))

	// Output:
	// FirstError() == nil: true
	// FirstError(x, y) == x
}
