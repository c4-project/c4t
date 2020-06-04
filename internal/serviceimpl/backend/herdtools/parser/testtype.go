// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package parser

import "fmt"

// TestType is the type of test we're parsing.
type TestType int

const (
	// None states that we haven't parsed a test type yet.
	None TestType = iota
	// Allowed is the 'allowed' test type.
	Allowed
	// Required is the 'required' test type.
	Required
)

// parseTestType tries to parse the test type from the word s.
func parseTestType(s string) (TestType, error) {
	switch s {
	case "Allowed":
		return Allowed, nil
	case "Required":
		return Required, nil
	default:
		return None, fmt.Errorf("%w: bad test type name %q", ErrBadTestType, s)
	}
}
