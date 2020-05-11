// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/helper/confhelp"
	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

// QuantitySet represents the part of a configuration that holds various tunable parameters for the batch runner.
type QuantitySet struct {
	// CorpusSize is the sampling size for the corpus after fuzzing.
	// It has a similar effect to CorpusSize in planner.Planner.
	CorpusSize int `toml:"corpus_size,omitzero"`

	// SubjectCycles is the number of times to fuzz each file.
	SubjectCycles int `toml:"subject_cycles,omitzero"`

	// NWorkers is the number of workers to use when fuzzing.
	NWorkers int `toml:"workers,omitzero"`
}

// Override substitutes any quantities in new that are non-zero for those in this set.
func (q *QuantitySet) Override(new QuantitySet) {
	confhelp.GenericOverride(q, new)
}

// Log logs q to l.
func (q *QuantitySet) Log(l *log.Logger) {
	confhelp.LogWorkers(l, q.NWorkers)
	l.Println("fuzzing each subject", iohelp.PluralQuantity(q.SubjectCycles, "time", "", "s"))
	l.Println("target corpus size:", iohelp.PluralQuantity(q.CorpusSize, "subject", "", "s"))
}
