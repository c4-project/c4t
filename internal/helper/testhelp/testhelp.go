// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package testhelp contains test helpers.
package testhelp

import (
	"errors"
	"testing"
)

// ExpectErrorIs checks whether got has an 'Is' relation to want (or, if want is nil, whether got is non-nil).
// If not, it fails the test with a message mentioning context and returns false.
func ExpectErrorIs(t *testing.T, got, want error, context string) bool {
	t.Helper()

	switch {
	case want == nil && got != nil:
		t.Errorf("%s: unexpected error: %q", context, got)
		return false
	case want != nil && got == nil:
		t.Errorf("%s: error nil; want=%q", context, want)
		return false
	case !errors.Is(got, want):
		t.Errorf("%s: error=%q; want=%q", context, got, want)
		return false
	}
	return true
}
