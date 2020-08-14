// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package quantity

import (
	"log"
)

// PlanSet contains configurable quantities for the planner.
type PlanSet struct {
	// NWorkers is the number of workers to use when probing the corpus.
	NWorkers int `toml:"workers,omitzero"`
}

// Override substitutes any quantities in new that are non-zero for those in this set.
func (q *PlanSet) Override(new PlanSet) {
	GenericOverride(q, new)
}

// Log logs q to l.
func (q *PlanSet) Log(l *log.Logger) {
	LogWorkers(l, q.NWorkers)
}
