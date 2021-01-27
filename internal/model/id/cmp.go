// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id

import "strings"

// Less compares two IDs lexicographically.
func (i ID) Less(i2 ID) bool {
	for j := 0; j < len(i.tags) && j < len(i2.tags); j++ {
		switch {
		case i.tags[j] < i2.tags[j]:
			return true
		case i.tags[j] > i2.tags[j]:
			return false
		}
	}
	return len(i.tags) < len(i2.tags)
}

// Equal compares two IDs for equality under case folding.
func (i ID) Equal(i2 ID) bool {
	li := len(i.tags)
	if li != len(i2.tags) {
		return false
	}
	for j := 0; j < li; j++ {
		if !strings.EqualFold(i.tags[j], i2.tags[j]) {
			return false
		}
	}
	return true
}
