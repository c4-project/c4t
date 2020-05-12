// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"fmt"
	"strings"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Analysis represents an analysis of a plan.
type Analysis struct {
	// ByStatus maps each status to the corpus of subjects that fall into it.
	ByStatus map[subject.Status]corpus.Corpus

	// Compilers maps each compiler ID to an analysis of that compiler.
	Compilers map[string]Compiler

	// Flags aggregates all flags found during the analysis.
	Flags Flag
}

// Compiler represents information about a compiler in a corpus analysis.
type Compiler struct {
	// Info contains the compiler's plan record.
	Info compiler.Compiler

	// Counts maps each status to the number of times it was observed across the corpus.
	Counts map[subject.Status]int

	Time *TimeSet
}

type TimeSet struct {
	Min  time.Duration
	Mean time.Duration
	Max  time.Duration
}

// String summarises this collation as a string.
func (a *Analysis) String() string {
	var sb strings.Builder

	bf := a.ByStatus

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
func (a *Analysis) HasFlagged() bool {
	return a.Flags.matches(FlagFlagged)
}

// HasFailures tests whether a collation has failure cases.
func (a *Analysis) HasFailures() bool {
	return a.Flags&(FlagCompileFail|FlagRunFail) != 0
}
