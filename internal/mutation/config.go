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
	Ranges []Range `toml:"ranges,omitempty"`
}

// Mutants returns a list of all mutant numbers to consider in this testing campaign.
//
// Mutants appear in the order defined, without deduplication.
func (c Config) Mutants() []uint {
	var m []uint

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

// Mutants expands a range into the slice of mutant numbers falling within it.
func (r Range) Mutants() []uint {
	if r.End <= r.Start {
		return []uint{}
	}

	m := make([]uint, r.End-r.Start)
	for i := r.Start; i < r.End; i++ {
		m[i-r.Start] = i
	}

	return m
}
