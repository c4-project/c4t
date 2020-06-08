// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"fmt"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/model/plan"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
)

// Analysis represents an analysis of a plan.
type Analysis struct {
	// Plan points to the plan that created this analysis.
	Plan *plan.Plan

	// ByStatus maps each status to the corpus of subjects that fall into it.
	ByStatus map[status.Status]corpus.Corpus

	// Compilers maps each compiler ID to an analysis of that compiler.
	Compilers map[string]Compiler

	// Flags aggregates all flags found during the analysis.
	Flags status.Flag
}

// Compiler represents information about a compiler in a corpus analysis.
type Compiler struct {
	// Info contains the compiler's plan record.
	Info compiler.Compiler

	// Counts maps each status to the number of times it was observed across the corpus.
	Counts map[status.Status]int

	// Time gathers statistics about how long, on average, this compiler took to compile corpus subjects.
	// It doesn't contain information about failed compilations.
	Time *TimeSet

	// RunTime gathers statistics about how long, on average, this compiler's compiled subjects took to run.
	// It doesn't contain information about failed compilations or runs (flagged runs are counted).
	RunTime *TimeSet
}

// String summarises this collation as a string.
func (a *Analysis) String() string {
	var sb strings.Builder

	bf := a.ByStatus

	// We range over this to enforce a deterministic order.
	for i := status.Ok; i <= status.Last; i++ {
		if i != status.Ok {
			sb.WriteString(", ")
		}
		_, _ = fmt.Fprintf(&sb, "%d %s", len(bf[i]), i.String())
	}

	return sb.String()
}

// HasFlagged tests whether a collation has flagged cases.
func (a *Analysis) HasFlagged() bool {
	return a.Flags.Matches(status.FlagFlagged)
}

// HasFailures tests whether a collation has failure cases.
func (a *Analysis) HasFailures() bool {
	return a.Flags&(status.FlagCompileFail|status.FlagRunFail) != 0
}
