// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate

import (
	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// ByStatus gets the corpi that make up this collation, mapped to each subject status.
func (c *Collation) ByStatus() map[subject.Status]corpus.Corpus {
	return map[subject.Status]corpus.Corpus{
		subject.StatusOk:             c.Successes,
		subject.StatusFlagged:        c.Flagged,
		subject.StatusCompileFail:    c.Compile.Failures,
		subject.StatusCompileTimeout: c.Compile.Timeouts,
		subject.StatusRunFail:        c.Run.Failures,
		subject.StatusRunTimeout:     c.Run.Timeouts,
	}
}

// HasFlagged tests whether a collation has flagged cases.
func (c *Collation) HasFlagged() bool {
	return len(c.Flagged) != 0
}

// HasFailures tests whether a collation has failure cases.
func (c *Collation) HasFailures() bool {
	return !(c.Compile.IsEmpty() && c.Run.IsEmpty())
}

// IsEmpty tests whether a fail collation has any results.
func (f *FailCollation) IsEmpty() bool {
	return len(f.Failures) == 0 && len(f.Timeouts) == 0
}
