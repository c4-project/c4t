// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package errhelp_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/helper/errhelp"
)

// ExampleFirstError is a testable example for FirstError.
func ExampleFirstError() {
	fmt.Println("FirstError() == nil:", errhelp.FirstError() == nil)

	fmt.Println("FirstError(x, y) ==", errhelp.FirstError(errors.New("x"), errors.New("y")))

	// Output:
	// FirstError() == nil: true
	// FirstError(x, y) == x
}

// ExampleTimeoutFirstError is a testable example for FirstError.
func ExampleTimeoutOrFirstError() {
	ctx := context.Background()
	fmt.Println("TimeoutOrFirstError(ok ctx) ==", errhelp.TimeoutOrFirstError(ctx))
	fmt.Println("TimeoutOrFirstError(ok ctx, x, y) ==", errhelp.TimeoutOrFirstError(ctx, errors.New("x"), errors.New("y")))

	tctx, _ := context.WithTimeout(ctx, 0)
	fmt.Println("TimeoutOrFirstError(t/o ctx) ==", errhelp.TimeoutOrFirstError(tctx))
	fmt.Println("TimeoutOrFirstError(t/o ctx, x, y) ==", errhelp.TimeoutOrFirstError(tctx, errors.New("x"), errors.New("y")))

	// Output:
	// TimeoutOrFirstError(ok ctx) == <nil>
	// TimeoutOrFirstError(ok ctx, x, y) == x
	// TimeoutOrFirstError(t/o ctx) == context deadline exceeded
	// TimeoutOrFirstError(t/o ctx, x, y) == context deadline exceeded
}
