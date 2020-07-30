// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package perturber

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/helper/confhelp"
	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

// QuantitySet contains configurable quantities for the planner.
type QuantitySet struct {
	// CorpusSize is the requested size of the test corpus.
	// If zero, no corpus sampling is done, but the planner will still error if the final corpus size is 0.
	// If nonzero, the corpus will be sampled if larger than the size, and an error occurs if the final size is below
	// that requested.
	CorpusSize int `toml:"corpus_size,omitzero"`
}

// Override substitutes any quantities in new that are non-zero for those in this set.
func (q *QuantitySet) Override(new QuantitySet) {
	confhelp.GenericOverride(q, new)
}

// Log logs q to l.
func (q *QuantitySet) Log(l *log.Logger) {
	l.Println("target corpus size:", iohelp.PluralQuantity(q.CorpusSize, "subject", "", "s"))
}
