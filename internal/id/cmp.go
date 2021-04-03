// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id

// Less compares two IDs lexicographically.
func (i ID) Less(i2 ID) bool {
	return lessTags(i.Tags(), i2.Tags())
}

func lessTags(t1 []string, t2 []string) bool {
	l1, l2 := len(t1), len(t2)
	for i := 0; i < l1 && i < l2; i++ {
		if t1[i] == t2[i] {
			continue
		}
		return t1[i] < t2[i]
	}
	return l1 < l2
}

// Equal compares two IDs for equality.
func (i ID) Equal(i2 ID) bool {
	// Using invariant that reprs are case-folded already.
	return i.repr == i2.repr
}
