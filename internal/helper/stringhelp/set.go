// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stringhelp

import "golang.org/x/exp/maps"

// Set is a string set, implemented in the usual Go way as a map to empty structs.
type Set map[string]struct{}

// NewSet constructs a set from the elements xs.
func NewSet(xs ...string) Set {
	s := make(Set, len(xs))
	s.Add(xs...)
	return s
}

// Add adds each element in xs to this set.
func (s Set) Add(xs ...string) {
	for _, x := range xs {
		s[x] = struct{}{}
	}
}

// Remove removes each element in xs from this set.
func (s Set) Remove(xs ...string) {
	for _, x := range xs {
		delete(s, x)
	}
}

// Copy makes a deep copy of this set.
func (s Set) Copy() Set {
	return maps.Clone(s)
}

// Slice returns the elements of xs as a sorted slice.
func (s Set) Slice() []string {
	return maps.Keys(s)
}
