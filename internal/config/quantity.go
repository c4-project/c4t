// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package config

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/controller/mach"
)

// QuantitySet is a set of tunable quantities for the director's stages.
type QuantitySet struct {
	// Fuzz is the quantity set for the fuzz stage.
	Fuzz fuzzer.QuantitySet `toml:"fuzz,omitzero"`
	// Mach is the quantity set for the machine-local stage, as well as any machine-local stages run remotely.
	Mach mach.QuantitySet `toml:"mach,omitzero"`
}

// Log logs q to l.
func (q *QuantitySet) Log(l *log.Logger) {
	l.Println("[Fuzzer]")
	q.Fuzz.Log(l)
	l.Println("[Mach]")
	q.Mach.Log(l)
}

// Override substitutes any quantities in new that are non-zero for those in this set.
func (q *QuantitySet) Override(new QuantitySet) {
	q.Fuzz.Override(new.Fuzz)
	q.Mach.Override(new.Mach)
}
