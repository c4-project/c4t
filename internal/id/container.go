// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id

import (
	"errors"
	"reflect"
	"sort"
)

// ErrNotMap occurs when we try to use an ID map function on something that isn't an ID map.
var ErrNotMap = errors.New("not a map with ID keys")

// Sort sorts ids.
func Sort(ids []ID) {
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].Less(ids[j])
	})
}

// MapKeys tries to get the keys of an ID-as-string map m as a sorted list.
// It fails if m is not an ID-as-string map.
func MapKeys(m interface{}) ([]ID, error) {
	ids, err := unsortedMapKeys(m)
	if err != nil {
		return nil, err
	}

	Sort(ids)
	return ids, nil
}

func unsortedMapKeys(m interface{}) ([]ID, error) {
	mv, _, err := checkMap(m)
	if err != nil {
		return nil, err
	}

	return unreflectIdSlice(mv.MapKeys()), nil
}

func unreflectIdSlice(kvs []reflect.Value) []ID {
	// Assuming we have already checked the type.
	ids := make([]ID, len(kvs))
	for i, kv := range kvs {
		ids[i] = kv.Interface().(ID)
	}
	return ids
}

// MapGlob filters a string map m to those keys that match glob when interpreted as IDs.
func MapGlob(m interface{}, glob ID) (interface{}, error) {
	mv, mt, err := checkMap(m)
	if err != nil {
		return nil, err
	}
	nm := reflect.MakeMap(mt)
	for _, kv := range mv.MapKeys() {
		match, err := kv.Interface().(ID).Matches(glob)
		if err != nil {
			return nil, err
		}
		if match {
			nm.SetMapIndex(kv, mv.MapIndex(kv))
		}
	}
	return nm.Interface(), nil
}

func checkMap(m interface{}) (reflect.Value, reflect.Type, error) {
	mv := reflect.ValueOf(m)
	if mv.Kind() != reflect.Map {
		return reflect.Value{}, nil, ErrNotMap
	}
	mt := mv.Type()
	if mt.Key() != reflect.TypeOf(ID{}) {
		return reflect.Value{}, nil, ErrNotMap
	}
	return mv, mt, nil
}

// LookupPrefix looks up id in map m by starting from the id itself, progressively taking a smaller and smaller
// prefix up to and including the empty ID, and returning the first value found.
// If a lookup succeeded, LookupPrefix returns the matched key as key, the value as val, and true as ok;
// else, it returns false, and the other two values are undefined.
func LookupPrefix(m interface{}, id ID) (key ID, val interface{}, ok bool) {
	mv, _, err := checkMap(m)
	if err != nil {
		return ID{}, nil, false
	}

	key = id
	for ok = true; ok; key, _, ok = key.Unsnoc() {
		vv := mv.MapIndex(reflect.ValueOf(key))
		if vv.Kind() != reflect.Invalid {
			return key, vv.Interface(), true
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
