// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
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

	// NWorkers is the number of workers to spawn in the worker pool.
	NWorkers int
}

// Override overrides this quantity set with the non-zero fields of other.
func (q *QuantitySet) Override(other QuantitySet) {
	if other.Count != 0 {
		q.Count = other.Count
	}
	if other.NWorkers != 0 {
		q.Count = other.NWorkers
	}
	if len(other.Divisions) != 0 {
		q.Divisions = other.Divisions
	}
}

// Bucket is the type of coverage buckets.
type Bucket struct {
	// Name is the name of the bucket, in the form "[0-9]+(,[0-9]+)*".
	Name string
	// Size is the size of the bucket.
	Size int
}

// String gets a string representation of this bucket.
func (b Bucket) String() string {
	return fmt.Sprintf("%s[%d]", b.Name, b.Size)
}

// Buckets calculates the set of buckets that should be constructed for the coverage setup.
// Buckets are allocated recursively according to the divisions set in the quantity set; each division carves the
// first bucket from the previous division into that many sub-buckets.
// Buckets always appear in reverse order, from highest outer bucket to lowest inner bucket.
func (q *QuantitySet) Buckets() []Bucket {
	var buckets []Bucket

	// This doesn't have particularly elegant properties, but it seems like it's the most obvious result.
	if q.Count <= 0 {
		return buckets
	}

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

		prefix := strings.Repeat("1_", i)
		for j := nbuckets; j >= 2; j-- {
			buckets = append(buckets, Bucket{Name: fmt.Sprintf("%s%d", prefix, j), Size: bsize})
		}

		bsize += brem
		if i == ndivs-1 {
			buckets = append(buckets, Bucket{Name: prefix + "1", Size: bsize})
		}
	}

	return buckets
}
