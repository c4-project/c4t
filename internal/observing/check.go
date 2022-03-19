// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package observing

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrObserverNil occurs when CheckObservers is passed a nil observer.
var ErrObserverNil = errors.New("observer nil")

// CheckObservers does some cursory checking of the observer slice obs (currently just nil checking).
func CheckObservers[O any](obs []O) error {
	for i, o := range obs {
		if !reflect.ValueOf(o).IsValid() {
			return fmt.Errorf("%w: index %d", ErrObserverNil, i)
		}
	}
	return nil
}
