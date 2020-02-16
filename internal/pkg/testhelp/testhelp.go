// Package testhelp contains test helpers.
package testhelp

import (
	"errors"
	"testing"
)

// ExpectErrorIs checks whether got has an 'Is' relation to want.
// If not, it fails the test with a message mentioning context.
func ExpectErrorIs(t *testing.T, got, want error, context string) {
	if got == nil {
		t.Helper()
		t.Errorf("%s: error nil; want=%q", context, want)
	} else if !errors.Is(got, want) {
		t.Helper()
		t.Errorf("%s: error=%q; want=%q", context, got, want)
	}
}
