// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package quantity

import (
	"log"

	"github.com/c4-project/c4t/internal/helper/stringhelp"
)

// PerturbSet contains configurable quantities for the perturber.
type PerturbSet struct {
	// CorpusSize is the requested size of the test corpus.
	// If zero, no corpus sampling is done, but the perturber will still error if the final corpus size is 0.
	// If nonzero, the corpus will be sampled if larger than the size, and an error occurs if the final size is below
	// that requested.
	CorpusSize int `toml:"corpus_size,omitzero" json:"corpus_size,omitempty"`
}

// Override substitutes any quantities in new that are non-zero for those in this set.
func (q *PerturbSet) Override(new PerturbSet) {
	GenericOverride(q, new)
}

// Log logs q to l.
func (q *PerturbSet) Log(l *log.Logger) {
	l.Println("target corpus size:", stringhelp.PluralQuantity(q.CorpusSize, "subject", "", "s"))
}
