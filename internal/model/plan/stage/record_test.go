// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stage_test

import (
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/plan/stage"
)

// ExampleNewRecord is a testable example for NewRecord.
func ExampleNewRecord() {
	start := time.Date(1997, time.May, 1, 23, 0, 0, 0, time.UTC)
	end := start.Add(10 * time.Minute)
	r := stage.NewRecord(stage.Fuzz, start, end)
	fmt.Println(r)

	// Output:
	// Fuzz completed on 1997-05-01T23:10:00Z (took 10m0s)
}
