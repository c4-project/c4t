// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package herdtools

import "fmt"

// TestType is the type of test we're parsing.
type TestType int

const (
	// TTNone states that we haven't parsed a test type yet.
	TTNone TestType = iota
	// TTAllowed is the 'allowed' test type.
	TTAllowed
	// TTRequired is the 'required' test type.
	TTRequired
)

// parseTestType tries to parse the test type from the word s.
func parseTestType(s string) (TestType, error) {
	switch s {
	case "Allowed":
		return TTAllowed, nil
	case "Required":
		return TTRequired, nil
	default:
		return TTNone, fmt.Errorf("%w: bad test type name %q", ErrBadTestType, s)
	}
}
