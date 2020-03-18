// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package id

import (
	"errors"
	"reflect"
	"sort"
)

var ErrNotMap = errors.New("not a map with string keys")

// Sort sorts ids.
func Sort(ids []ID) {
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].Less(ids[j])
	})
}

// MapKeys tries to get the keys of an ID-as-string map m as a sorted list.
// It fails if m is not an ID-as-string map.
func MapKeys(m interface{}) ([]ID, error) {
	mv := reflect.ValueOf(m)
	mt := mv.Type()
	if mt.Kind() != reflect.Map || mt.Key().Kind() != reflect.String {
		return nil, ErrNotMap
	}

	keys := mv.MapKeys()
	ids := make([]ID, len(keys))
	for i := range keys {
		var err error
		if ids[i], err = TryFromString(keys[i].String()); err != nil {
			return nil, err
		}
	}

	Sort(ids)
	return ids, nil
}
