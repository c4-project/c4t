// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import "time"

// TimeSet contains aggregate statistics about timing (of compilers, runs, etc).
type TimeSet struct {
	// Min contains the minimum time reported across the corpus.
	Min time.Duration
	// Sum contains the total time reported across the corpus.
	Sum time.Duration
	// Max contains the maximum time reported across the corpus.
	Max time.Duration
	// Count contains the number of samples taken for this time set.
	Count int
}

// NewTimeSet produces a timeset from the given raw times.
func NewTimeSet(raw ...time.Duration) *TimeSet {
	var t TimeSet

	for _, r := range raw {
		t.add(r)
	}

	return &t
}

// add logs r if it is the minimum or maximum time, and adds it to the mean.
// Note that this does not calculate a rolling mean, but instead a sum; the .Mean field will need to be divided
// once all adds are done.
func (t *TimeSet) add(r time.Duration) {
	t.Sum += r
	t.Count++

	if t.Min == 0 || r < t.Min {
		t.Min = r
	}
	if t.Max == 0 || t.Max < r {
		t.Max = r
	}
}

// Mean calculates the arithmetic mean of the times in this set.
func (t TimeSet) Mean() time.Duration {
	if t.Count == 0 {
		return 0
	}
	return t.Sum / time.Duration(t.Count)
}
