// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/helper/confhelp"
)

// QuantitySet contains configurable quantities for the planner.
type QuantitySet struct {
	// NWorkers is the number of workers to use when probing the corpus.
	NWorkers int `toml:"workers,omitzero"`
}

// Override substitutes any quantities in new that are non-zero for those in this set.
func (q *QuantitySet) Override(new QuantitySet) {
	confhelp.GenericOverride(q, new)
}

// Log logs q to l.
func (q *QuantitySet) Log(l *log.Logger) {
	confhelp.LogWorkers(l, q.NWorkers)
}
