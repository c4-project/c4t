// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gccnt

import "sort"

// ConditionSet contains the sets of conditions at which gccn't will fail.
type ConditionSet struct {
	// Diverge contains the conditions at which gccn't will diverge.
	Diverge Condition

	// Error contains the conditions at which gccn't will error.
	Error Condition

	// MutHitPeriod is an integer that, when nonzero, is the period in the mutant number that will trigger gccn't
	// reporting a mutation hit without an error.
	MutHitPeriod uint64
}

// sort makes sure the opt lists in the conditionset are sorted.
func (c *ConditionSet) sort() {
	c.Diverge.sort()
	c.Error.sort()
}

// Condition specifies a condition at which gccn't will fail in some specified way.
type Condition struct {
	// Opts is a list of optimisation levels that will trigger this error.
	Opts []string
	// MutPeriod is an integer that, when nonzero, is the period in the mutant number that will trigger this error.
	MutPeriod uint64
}

// sort makes sure the opt lists in the condition are sorted.
func (c *Condition) sort() {
	sort.Strings(c.Opts)
}
