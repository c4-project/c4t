// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package id

import (
	"reflect"
	"sort"

	"github.com/MattWindsor91/act-tester/internal/helper/stringhelp"
)

// Sort sorts ids.
func Sort(ids []ID) {
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].Less(ids[j])
	})
}

// MapKeys tries to get the keys of an ID-as-string map m as a sorted list.
// It fails if m is not an ID-as-string map.
func MapKeys(m interface{}) ([]ID, error) {
	keys, err := stringhelp.SafeMapKeys(m)
	if err != nil {
		return nil, err
	}

	ids := make([]ID, len(keys))
	for i := range keys {
		var err error
		if ids[i], err = tryFromValue(keys[i]); err != nil {
			return nil, err
		}
	}

	Sort(ids)
	return ids, nil
}

// MapGlob filters a string map m to those keys that match glob when interpreted as IDs.
func MapGlob(m interface{}, glob ID) (interface{}, error) {
	mv, mt, err := stringhelp.CheckMap(m)
	if err != nil {
		return nil, err
	}

	nm := reflect.MakeMap(mt)
	for _, kstr := range mv.MapKeys() {
		k, err := tryFromValue(kstr)
		if err != nil {
			return nil, err
		}
		match, err := k.Matches(glob)
		if err != nil {
			return nil, err
		}
		if match {
			nm.SetMapIndex(kstr, mv.MapIndex(kstr))
		}
	}
	return nm.Interface(), nil
}

func tryFromValue(v reflect.Value) (ID, error) {
	return TryFromString(v.String())
}
