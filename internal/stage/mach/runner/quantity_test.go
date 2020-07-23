// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/runner"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/timeout"
)

// ExampleQuantitySet_Log is a runnable example for Log.
func ExampleQuantitySet_Log() {
	l := log.New(os.Stdout, "", 0)

	fmt.Println("[empty]")
	var qs runner.QuantitySet
	qs.Log(l)

	fmt.Println("[with 5 workers]")
	qs.NWorkers = 5
	qs.Log(l)

	fmt.Println("[and a 1 minute timeout]")
	qs.Timeout = timeout.Timeout(1 * time.Minute)
	qs.Log(l)

	// Output:
	// [empty]
	// running across 0 workers
	// [with 5 workers]
	// running across 5 workers
	// [and a 1 minute timeout]
	// running across 5 workers
	// timeout at 1m0s
}
