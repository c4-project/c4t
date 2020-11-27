// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package observing

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrNotSlice occurs when CheckObservers is passed something other than a slice.
	ErrNotSlice = errors.New("not an observer slice")

	// ErrObserverNil occurs when CheckObservers is passed a nil observer.
	ErrObserverNil = errors.New("observer nil")
)

// CheckObservers does some cursory checking of the observer slice obs (currently just nil checking).
func CheckObservers(obs interface{}) error {
	rv := reflect.ValueOf(obs)
	if rv.Kind() != reflect.Slice {
		return fmt.Errorf("%w: %v", ErrNotSlice, obs)
	}
	l := rv.Len()
	for i := 0; i < l; i++ {
		if rv.Index(i).IsNil() {
			return fmt.Errorf("%w: index %d", ErrObserverNil, i)
		}
	}
	return nil
}
