// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package quantity

import (
	"log"
)

// MachNodeSet contains the tunable quantities for both batch-compiler and batch-runner.
type MachNodeSet struct {
	// Compiler is the quantity set for the compiler.
	Compiler BatchSet `toml:"compiler,omitzero" json:"compiler,omitempty"`
	// Runner is the quantity set for the runner.
	Runner BatchSet `toml:"runner,omitzero" json:"runner,omitempty"`
}

// Log logs q to l.
func (q *MachNodeSet) Log(l *log.Logger) {
	l.Println("[Compiler]")
	q.Compiler.Log(l)
	l.Println("[Runner]")
	q.Runner.Log(l)
}

// Override overrides the quantities in this set with any new quantities supplied in new.
func (q *MachNodeSet) Override(new MachNodeSet) {
	q.Compiler.Override(new.Compiler)
	q.Runner.Override(new.Runner)
}

// BatchSet contains the tunable quantities for either a batch compiler or a batch runner.
type BatchSet struct {
	// Timeout is the timeout for each runner.
	// Non-positive values disable the timeout.
	Timeout Timeout `toml:"timeout,omitzero" json:"timeout,omitempty"`

	// NWorkers is the number of parallel run workers that should be spawned.
	// Anything less than or equal to 1 will sequentialise the run.
	NWorkers int `toml:"workers,omitzero" json:"workers,omitempty"`
}

// Log logs this quantity set to l.
func (q *BatchSet) Log(l *log.Logger) {
	LogWorkers(l, q.NWorkers)
	q.Timeout.Log(l)
}

// Override substitutes any non-zero quantities in new for those in this quantity set, in-place.
func (q *BatchSet) Override(new BatchSet) {
	if new.Timeout.IsActive() {
		q.Timeout = new.Timeout
	}
	if new.NWorkers != 0 {
		q.NWorkers = new.NWorkers
	}
}
