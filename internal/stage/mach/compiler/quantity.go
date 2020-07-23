// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/timeout"
)

// QuantitySet contains tunable quantities for the batch-compiler.
type QuantitySet struct {
	// Timeout is the timeout for each compile.
	// Non-positive values disable the timeout.
	Timeout timeout.Timeout `toml:"timeout,omitzero"`
}

func (q *QuantitySet) Log(l *log.Logger) {
	q.Timeout.Log(l)
}

// Override substitutes any non-zero quantities in new for those in this quantity set, in-place.
func (q *QuantitySet) Override(new QuantitySet) {
	if new.Timeout.IsActive() {
		q.Timeout = new.Timeout
	}
}
