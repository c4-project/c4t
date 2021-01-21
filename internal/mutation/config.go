// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation

// Config configures a particular mutation testing campaign.
//
// This currently just tracks ranges of mutation numbers, but may be generalised if we branch to supporting more than
// one kind of mutation test.
type Config struct {
	// Enabled gets whether mutation testing is enabled.
	//
	// Setting this to false is equivalent to setting Ranges to empty.
	Enabled bool `toml:"enabled,omitempty"`

	// Ranges contains the list of mutation number ranges that the campaign should use.
	Ranges []Range `toml:"ranges,omitempty"`
}

// IsActive gets whether this Config is enabled and has a functional set of ranges, without evaluating the mutant set.
func (c Config) IsActive() bool {
	if !c.Enabled {
		return false
	}
	for _, r := range c.Ranges {
		if !r.IsEmpty() {
			return true
		}
	}
	return false
}

// Mutants returns a list of all mutant numbers to consider in this testing campaign.
//
// Mutants appear in the order defined, without deduplication.
// If Enabled is false, Mutants will be empty.
func (c Config) Mutants() []uint {
	var m []uint

	if !c.Enabled {
		return m
	}

	for _, r := range c.Ranges {
		m = append(m, r.Mutants()...)
	}

	return m
}

// Range defines an inclusive numeric range of mutant numbers to consider.
type Range struct {
	// Start is the first mutant number to consider in this range.
	Start uint `toml:"from"`
	// End is one past the last mutant number to consider in this range.
	End uint `toml:"to"`
}

// IsEmpty gets whether this range defines no mutant numbers.
func (r Range) IsEmpty() bool {
	return r.End <= r.Start
}

// Mutants expands a range into the slice of mutant numbers falling within it.
func (r Range) Mutants() []uint {
	if r.IsEmpty() {
		return []uint{}
	}

	m := make([]uint, r.End-r.Start)
	for i := r.Start; i < r.End; i++ {
		m[i-r.Start] = i
	}

	return m
}
