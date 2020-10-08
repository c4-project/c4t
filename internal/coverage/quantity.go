// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"fmt"
	"strings"
)

// QuantitySet contains the quantities tracked by the coverage generator.
type QuantitySet struct {
	// Count is the number of subjects to fuzz for each profile.
	Count int `toml:"count"`

	// Divisions specifies how to divide Count subjects into buckets.
	// Divisions behave recursively: each subsequent level of division gets applied to the first bucket in the
	// previous level.
	Divisions []int `toml:"divisions"`
}

// Override overrides this quantity set with the non-zero fields of other.
func (q *QuantitySet) Override(other QuantitySet) {
	if other.Count != 0 {
		q.Count = other.Count
	}
	if len(other.Divisions) != 0 {
		q.Divisions = other.Divisions
	}
}

// Buckets calculates the set of buckets that should be constructed for the coverage setup.
// Buckets are allocated recursively according to the divisions set in the quantity set; each division carves the
// first bucket from the previous division into that many sub-buckets.
func (q *QuantitySet) Buckets() map[string]int {
	buckets := map[string]int{}

	// No divisions should be the same as a single all-encompassing division.
	divs := q.Divisions
	if len(divs) == 0 {
		divs = []int{1}
	}

	ndivs := len(divs)
	bsize := q.Count
	for i, nbuckets := range divs {
		// Avoid divisions by zero, etc.
		if nbuckets <= 0 {
			nbuckets = 1
		}

		// If we don't divide cleanly, distribute the remainder to bucket 1
		brem := bsize % nbuckets
		bsize /= nbuckets

		prefix := strings.Repeat("1,", i)
		for j := 2; j <= nbuckets; j++ {
			buckets[fmt.Sprintf("%s%d", prefix, j)] = bsize
		}

		bsize += brem
		if i == ndivs-1 {
			buckets[prefix+"1"] = bsize
		}
	}

	return buckets
}
