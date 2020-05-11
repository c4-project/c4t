// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"fmt"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Analysis represents an analysis of a corpus.
type Analysis struct {
	// ByStatus maps each status to the corpus of subjects that fall into it.
	ByStatus map[subject.Status]corpus.Corpus
}

// String summarises this collation as a string.
func (c *Analysis) String() string {
	var sb strings.Builder

	bf := c.ByStatus

	// We range over this to enforce a deterministic order.
	for i := subject.StatusOk; i < subject.NumStatus; i++ {
		if i != subject.StatusOk {
			sb.WriteString(", ")
		}
		_, _ = fmt.Fprintf(&sb, "%d %s", len(bf[i]), i.String())
	}

	return sb.String()
}

// HasFlagged tests whether a collation has flagged cases.
func (c *Analysis) HasFlagged() bool {
	return len(c.ByStatus[subject.StatusFlagged]) != 0
}

// HasFailures tests whether a collation has failure cases.
func (c *Analysis) HasFailures() bool {
	for i := subject.FirstBadStatus; i < subject.NumStatus; i++ {
		if len(c.ByStatus[i]) != 0 {
			return true
		}
	}
	return false
}
