// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package timing_test

import (
	"fmt"
	"time"

	"github.com/c4-project/c4t/internal/timing"
)

// ExampleSpanFromDuration is a testable example for SpanFromDuration.
func ExampleSpanFromDuration() {
	ts := timing.SpanFromDuration(time.Date(1990, time.January, 1, 12, 00, 00, 00, time.UTC), 10*time.Minute)
	fmt.Println(ts)
	fmt.Println(ts.IsInstant())
	fmt.Println(ts.IsUndefined())
	fmt.Printf("%.0f", ts.Duration().Minutes())

	// Output:
	// 10m0s (from 1990-01-01T12:00:00Z to 1990-01-01T12:10:00Z)
	// false
	// false
	// 10
}

// ExampleSpanFromInstant is a testable example for SpanFromInstant.
func ExampleSpanFromInstant() {
	ts := timing.SpanFromInstant(time.Date(1990, time.January, 1, 12, 00, 00, 00, time.UTC))
	fmt.Println(ts)
	fmt.Println(ts.IsInstant())
	fmt.Println(ts.IsUndefined())
	fmt.Printf("%.0f", ts.Duration().Minutes())

	// Output:
	// 1990-01-01T12:00:00Z
	// true
	// false
	// 0
}
