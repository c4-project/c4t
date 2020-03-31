// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package optlevel

import (
	"errors"
	"fmt"
)

// ErrNoSuchLevel occurs if a Selection enables an optimisation level that isn't available.
var ErrNoSuchLevel = errors.New("no such optimisation level")

// Resolver is the interface of types that support optimisation level lookup.
type Resolver interface {
	// DefaultLevels retrieves a set of optimisation levels that are enabled by default.
	DefaultLevels() map[string]struct{}
	// Levels retrieves a set of potential optimisation levels.
	// This map shouldn't be modified, as it may be global.
	Levels() map[string]Level
}

// SelectLevels selects from r the optimisation levels mentioned by s.
func SelectLevels(r Resolver, s Selection) (map[string]Level, error) {
	all := r.Levels()

	choices := s.Override(r.DefaultLevels())
	chosen := make(map[string]Level, len(choices))

	for c := range choices {
		var ok bool
		if chosen[c], ok = all[c]; !ok {
			return nil, fmt.Errorf("%w: %q", ErrNoSuchLevel, c)
		}
	}

	return chosen, nil
}
