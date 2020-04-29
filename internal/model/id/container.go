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

// LookupPrefix looks up id in map m by starting from the id itself, and progressively taking a smaller and smaller
// prefix until the id runs out or m has a value.
// If a lookup succeeded, LookupPrefix returns the matched key as key, the value as val, and true as ok;
// else, it returns false, and the other two values are undefined.
func LookupPrefix(m interface{}, id ID) (key ID, val interface{}, ok bool) {
	mv, _, err := stringhelp.CheckMap(m)
	if err != nil {
		return ID{}, nil, false
	}

	key = id
	for ok = true; ok; key, _, ok = key.Unsnoc() {
		vv := mv.MapIndex(reflect.ValueOf(key.String()))
		if vv.Kind() != reflect.Invalid {
			return key, vv.Interface(), true
		}
	}
	return ID{}, nil, ok
}

func tryFromValue(v reflect.Value) (ID, error) {
	return TryFromString(v.String())
}
