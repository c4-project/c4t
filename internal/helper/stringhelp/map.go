// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stringhelp

import (
	"errors"
	"reflect"
)

// ErrNotMap occurs when a generic string map function receives something that isn't a string map.
var ErrNotMap = errors.New("not a map with string keys")

// CheckMap checks, through reflection, to see if m is a map with string keys.
func CheckMap(m interface{}) (reflect.Value, reflect.Type, error) {
	mv := reflect.ValueOf(m)
	mt := mv.Type()
	if mt.Kind() != reflect.Map || mt.Key().Kind() != reflect.String {
		return reflect.Value{}, nil, ErrNotMap
	}
	return mv, mt, nil
}

// SafeMapKeys is the same as calling MapKeys on m, but checks to make sure m is a string map first.
func SafeMapKeys(m interface{}) ([]reflect.Value, error) {
	v, _, err := CheckMap(m)
	if err != nil {
		return nil, err
	}
	return v.MapKeys(), nil
}
