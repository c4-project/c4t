// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id

import (
	"golang.org/x/exp/maps"
	"sort"
)

// Sort sorts ids.
func Sort(ids []ID) {
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].Less(ids[j])
	})
}

// MapKeys gets the keys of an ID-as-string map m as a sorted list.
func MapKeys[V any](m map[ID]V) []ID {
	ids := maps.Keys(m)
	Sort(ids)
	return ids
}

// MapGlob filters a string map m to those keys that match glob when interpreted as IDs.
func MapGlob[V any](m map[ID]V, glob ID) (map[ID]V, error) {
	nm := make(map[ID]V, len(m))
	for k, v := range m {
		match, err := k.Matches(glob)
		if err != nil {
			return nil, err
		}
		if match {
			nm[k] = v
		}
	}
	return nm, nil
}

// LookupPrefix looks up id in map m by starting from the id itself, progressively taking a smaller and smaller
// prefix up to and including the empty ID, and returning the first value found.
// If a lookup succeeded, LookupPrefix returns the matched key as key, the value as val, and true as ok;
// else, it returns false, and the other two values are undefined.
func LookupPrefix[V any](m map[ID]V, id ID) (key ID, val interface{}, ok bool) {
	key = id
	for ok = true; ok; key, _, ok = key.Unsnoc() {
		if v, vok := m[key]; vok {
			return key, v, true
		}
	}
	return ID{}, nil, ok
}

// SearchSlice finds the smallest index in haystack for which the ID at that index is greater than or equal to needle.
// If there is no such index, it returns len(haystack).
func SearchSlice(haystack []ID, needle ID) int {
	return sort.Search(len(haystack), func(i int) bool {
		return !haystack[i].Less(needle)
	})
}
