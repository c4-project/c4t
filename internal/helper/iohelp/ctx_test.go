// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

// ExampleCheckDone is a runnable example for CheckDone.
func ExampleCheckDone() {
	ctx, cancel := context.WithCancel(context.Background())

	// This shouldn't block, but should return without an error:
	err1 := iohelp.CheckDone(ctx)

	// This should also not block, but return a 'cancelled' error:
	cancel()
	err2 := iohelp.CheckDone(ctx)

	fmt.Printf("1: nil=%v, cancelled=%v\n", err1 == nil, errors.Is(err1, context.Canceled))
	fmt.Printf("2: nil=%v, cancelled=%v\n", err2 == nil, errors.Is(err2, context.Canceled))

	// Output:
	// 1: nil=true, cancelled=false
	// 2: nil=false, cancelled=true
}
