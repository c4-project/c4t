// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package quantity

import (
	"log"

	"github.com/c4-project/c4t/internal/helper/stringhelp"
)

// FuzzSet represents the part of a configuration that holds various tunable parameters for the fuzzer.
type FuzzSet struct {
	// CorpusSize is the sampling size for the corpus after fuzzing.
	// It has a similar effect to CorpusSize in planner.Planner.
	CorpusSize int `toml:"corpus_size,omitzero" json:"corpus_size,omitempty"`

	// SubjectCycles is the number of times to fuzz each file.
	SubjectCycles int `toml:"subject_cycles,omitzero" json:"subject_cycles,omitempty"`

	// NWorkers is the number of workers to use when fuzzing.
	NWorkers int `toml:"workers,omitzero" json:"num_workers,omitempty"`
}

// Override substitutes any quantities in new that are non-zero for those in this set.
func (q *FuzzSet) Override(new FuzzSet) {
	GenericOverride(q, new)
}

// Log logs q to l.
func (q *FuzzSet) Log(l *log.Logger) {
	LogWorkers(l, q.NWorkers)
	l.Println("fuzzing each subject", stringhelp.PluralQuantity(q.SubjectCycles, "time", "", "s"))
	l.Println("target corpus size:", stringhelp.PluralQuantity(q.CorpusSize, "subject", "", "s"))
}
