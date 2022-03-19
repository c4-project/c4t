// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id

import "golang.org/x/exp/slices"

// Less compares two IDs lexicographically.
func (i ID) Less(i2 ID) bool {
	return slices.Compare(i.Tags(), i2.Tags()) < 0
}

// Equal compares two IDs for equality.
func (i ID) Equal(i2 ID) bool {
	// Using invariant that reprs are case-folded already.
	return i.repr == i2.repr
}
