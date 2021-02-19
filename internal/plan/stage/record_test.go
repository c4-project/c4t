// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stage_test

import (
	"fmt"
	"time"

	"github.com/c4-project/c4t/internal/timing"

	"github.com/c4-project/c4t/internal/plan/stage"
)

// ExampleRecord_String is a testable example for Record.String.
func ExampleRecord_String() {
	fmt.Println(
		stage.Record{
			Stage:    stage.Fuzz,
			Timespan: timing.SpanFromDuration(timing.MockDate, 10*time.Minute),
		},
	)

	// Output:
	// Fuzz: 10m0s (from 1997-05-01T21:00:00Z to 1997-05-01T21:10:00Z)
}
