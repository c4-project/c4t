// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation

import (
	"github.com/c4-project/c4t/internal/quantity"
)

// Config configures a particular mutation testing campaign.
//
// This currently just tracks ranges of mutation numbers, but may be generalised if we branch to supporting more than
// one kind of mutation test.
type Config struct {
	// Enabled gets whether mutation testing is enabled.
	//
	// Setting this to false is equivalent to setting Ranges to empty.
	Enabled bool `json:"enabled,omitempty" toml:"enabled,omitempty"`

	// Selection contains any selected mutation.
	//
	// This can be set in the tester's config file, in addition to or instead of Ranges, but will be overridden by
	// any automatic mutant selection.
	Selection Mutant `json:"selection,omitempty" toml:"selection,omitempty"`

	// Auto gathers configuration about how to automate mutation selection.
	Auto AutoConfig `json:"auto,omitempty" toml:"auto,omitempty"`
}

// AutoConfig specifies configuration pertaining to automatically selecting mutants.
//
// The mutation tester can be used with a manual selection, but is probably not very exciting.
type AutoConfig struct {
	// Ranges contains the list of mutation number ranges that the campaign should use for automatic mutant selection.
	Ranges []Range `json:"ranges,omitempty" toml:"ranges,omitempty"`

	// ChangeMutantAfter is the (minimum) duration that each mutant gets before being automatically incremented.
	// If 0, this sort of auto-increment
	ChangeAfter quantity.Timeout `json:"change_after,omitempty" toml:"change_after,omitempty"`

	// ChangeKilled specifies whether mutants should be automatically incremented after being killed.
	ChangeKilled bool `json:"change_killed" toml:"change_killed"`
}

// IsActive gets whether automatic selection is enabled.
func (c AutoConfig) IsActive() bool {
	return c.HasRanges() && (c.HasChangeAfter() || c.ChangeKilled)
}

// HasRanges gets whether at least one viable range exists, without expanding the ranges themselves.
func (c AutoConfig) HasRanges() bool {
	for _, r := range c.Ranges {
		if !r.IsEmpty() {
			return true
		}
	}
	return false
}

// HasChangeAfter gets whether ChangeAfter is set to something other than zero.
func (c AutoConfig) HasChangeAfter() bool {
	return c.ChangeAfter.IsActive()
}

// Mutants returns a list of all mutant numbers to consider in this testing campaign.
//
// Mutants appear in the order defined, without deduplication.
// If Enabled is false, Mutants will be empty.
func (c AutoConfig) Mutants() []Mutant {
	var m []Mutant

	for _, r := range c.Ranges {
		m = append(m, r.Mutants()...)
	}

	return m
}

// Range defines an inclusive numeric range of mutant numbers to consider.
type Range struct {
	// Operator is, if given, the name of the operator in this range.
	Operator string `json:"operator" toml:"operator"`

	// Start is the first mutant number to consider in this range.
	Start Index `json:"start" toml:"start"`
	// End is one past the last mutant number to consider in this range.
	End Index `json:"end" toml:"end"`
}

// IsEmpty gets whether this range defines no mutant numbers.
func (r Range) IsEmpty() bool {
	return r.End <= r.Start
}

// IsSingleton gets whether this range has exactly one item in it.
func (r Range) IsSingleton() bool {
	return r.End == r.Start+1
}

// Mutants expands a range into the slice of mutant numbers falling within it.
func (r Range) Mutants() []Mutant {
	switch {
	case r.IsEmpty():
		return []Mutant{}
	case r.IsSingleton():
		// Don't record variant numbers for singleton mutants
		return []Mutant{NamedMutant(r.Start, r.Operator, 0)}
	default:
		return r.enumMutants()
	}
}

func (r Range) enumMutants() []Mutant {
	m := make([]Mutant, r.End-r.Start)
	for i := r.Start; i < r.End; i++ {
		j := uint64(i - r.Start)
		m[j] = NamedMutant(i, r.Operator, j+1)
	}
	return m
}
