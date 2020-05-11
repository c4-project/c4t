// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate

import (
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// HasFlagged tests whether a collation has flagged cases.
func (c *Collation) HasFlagged() bool {
	return len(c.ByStatus[subject.StatusFlagged]) != 0
}

// HasFailures tests whether a collation has failure cases.
func (c *Collation) HasFailures() bool {
	for i := subject.FirstBadStatus; i < subject.NumStatus; i++ {
		if len(c.ByStatus[i]) != 0 {
			return true
		}
	}
	return false
}
