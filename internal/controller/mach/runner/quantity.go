// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/controller/mach/timeout"
)

// QuantitySet contains tunable quantities for the batch-runner.
type QuantitySet struct {
	// Timeout is the timeout for each runner.
	// Non-positive values disable the timeout.
	Timeout timeout.Timeout `toml:"timeout,omitzero"`

	// NWorkers is the number of parallel run workers that should be spawned.
	// Anything less than or equal to 1 will sequentialise the run.
	NWorkers int `toml:"workers,omitzero"`
}

// Log logs this quantity set to l.
func (q *QuantitySet) Log(l *log.Logger) {
	l.Println("running across", iohelp.PluralQuantity(q.NWorkers, "worker", "", "s"))
	q.Timeout.Log(l)
}
