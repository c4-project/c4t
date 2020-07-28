// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package quantity contains the quantity sets for the machine node.
package quantity

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/helper/confhelp"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/timeout"
)

// SingleSet contains the tunable quantities for either a batch compiler or a batch runner.
type SingleSet struct {
	// Timeout is the timeout for each runner.
	// Non-positive values disable the timeout.
	Timeout timeout.Timeout `toml:"timeout,omitzero"`

	// NWorkers is the number of parallel run workers that should be spawned.
	// Anything less than or equal to 1 will sequentialise the run.
	NWorkers int `toml:"workers,omitzero"`
}

// Log logs this quantity set to l.
func (q *SingleSet) Log(l *log.Logger) {
	confhelp.LogWorkers(l, q.NWorkers)
	q.Timeout.Log(l)
}

// Override substitutes any non-zero quantities in new for those in this quantity set, in-place.
func (q *SingleSet) Override(new SingleSet) {
	if new.Timeout.IsActive() {
		q.Timeout = new.Timeout
	}
	if new.NWorkers != 0 {
		q.NWorkers = new.NWorkers
	}
}

// Set contains the tunable quantities for both batch-compiler and batch-runner.
type Set struct {
	// Compiler is the quantity set for the compiler.
	Compiler SingleSet `toml:"compiler,omitzero"`
	// Runner is the quantity set for the runner.
	Runner SingleSet `toml:"runner,omitzero"`
}

// Log logs q to l.
func (q *Set) Log(l *log.Logger) {
	l.Println("[Compiler]")
	q.Compiler.Log(l)
	l.Println("[Runner]")
	q.Runner.Log(l)
}

// Override overrides the quantities in this set with any new quantities supplied in new.
func (q *Set) Override(new Set) {
	q.Compiler.Override(new.Compiler)
	q.Runner.Override(new.Runner)
}
