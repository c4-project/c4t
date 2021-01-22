// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"fmt"
	"strings"

	"github.com/c4-project/c4t/internal/mutation"

	"github.com/c4-project/c4t/internal/plan"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/c4-project/c4t/internal/subject/corpus"
)

// Analysis represents an analysis of a plan.
type Analysis struct {
	// Plan points to the plan that created this analyser.
	Plan *plan.Plan

	// ByStatus maps each status to the corpus of subjects that fall into it.
	ByStatus map[status.Status]corpus.Corpus

	// Compilers maps each compiler ID (or full-ID, depending on configuration) to an analysis of that compiler.
	Compilers map[string]Compiler

	// Flags aggregates all flags found during the analysis.
	Flags status.Flag

	// Mutation, if non-nil, contains information about mutation testing done over this plan.
	Mutation mutation.Analysis
}

// Compiler represents information about a compiler in a corpus analysis.
type Compiler struct {
	// Info contains the compiler's plan record.
	Info compiler.Instance

	// Counts maps each status to the number of times it was observed across the corpus.
	Counts map[status.Status]int

	// Logs maps each subject name to its compiler log.
	Logs map[string]string

	// Time gathers statistics about how long, on average, this compiler took to compile corpus subjects.
	// It doesn't contain information about failed compilations.
	Time *TimeSet

	// RunTime gathers statistics about how long, on average, this compiler's compiled subjects took to run.
	// It doesn't contain information about failed compilations or runs (flagged runs are counted).
	RunTime *TimeSet
}

func newAnalysis(p *plan.Plan) *Analysis {
	return &Analysis{
		Plan:      p,
		ByStatus:  make(map[status.Status]corpus.Corpus, status.Last),
		Compilers: make(map[string]Compiler, len(p.Compilers)),
		Mutation:  make(mutation.Analysis),
	}
}

// String summarises this collation as a string.
func (a *Analysis) String() string {
	var sb strings.Builder

	bf := a.ByStatus
	first := true

	// We range over this to enforce a deterministic order.
	for i := status.Ok; i <= status.Last; i++ {
		l := len(bf[i])
		if l == 0 {
			continue
		}
		if !first {
			sb.WriteString(", ")
		}
		first = false
		_, _ = fmt.Fprintf(&sb, "%d %s", l, i.String())
	}

	return sb.String()
}

// HasFlagged tests whether an analysis has flagged cases.
func (a *Analysis) HasFlagged() bool {
	return a.Flags.MatchesStatus(status.Flagged)
}

// HasFailures tests whether a collation has failure cases.
func (a *Analysis) HasFailures() bool {
	return a.Flags&(status.FlagFail) != 0
}

// HasBadOutcomes tests whether a collation has any bad (flagged, failed, or timed-out) cases.
func (a *Analysis) HasBadOutcomes() bool {
	return a.Flags&(status.FlagBad) != 0
}
