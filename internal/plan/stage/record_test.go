// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stage_test

import (
	"fmt"
	"time"

	"github.com/c4-project/c4t/internal/plan/stage"
)

// ExampleNewRecord is a testable example for NewRecord.
func ExampleNewRecord() {
	start := time.Date(1997, time.May, 1, 23, 0, 0, 0, time.UTC)
	r := stage.NewRecord(stage.Fuzz, start, 10*time.Minute)
	fmt.Println(r)

	// Output:
	// Fuzz completed on 1997-05-01T23:10:00Z (took 10m0s)
}
