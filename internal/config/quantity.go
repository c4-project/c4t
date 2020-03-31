// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package config

import "github.com/MattWindsor91/act-tester/internal/controller/fuzzer"

// QuantitySet is a set of tunable quantities for the director's stages.
type QuantitySet struct {
	// Fuzz is the quantity set for the fuzz stage.
	Fuzz fuzzer.QuantitySet `toml:"fuzz"`
}

// Override substitutes any quantities in new that are non-zero for those in this set.
func (q *QuantitySet) Override(new QuantitySet) {
	q.Fuzz.Override(new.Fuzz)
}
