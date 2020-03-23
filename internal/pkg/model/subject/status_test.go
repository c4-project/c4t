// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject_test

import (
	"strconv"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/helpers/testhelp"
)

// TestStatusOfString_RoundTrip checks that converting a status to and back from its string is the identity.
func TestStatusOfString_RoundTrip(t *testing.T) {
	t.Parallel()
	for want := subject.StatusUnknown; want < subject.NumStatus; want++ {
		want := want
		t.Run(strconv.Itoa(int(want)), func(t *testing.T) {
			t.Parallel()
			got, err := subject.StatusOfString(want.String())
			if err != nil {
				t.Errorf("unexpected error round-tripping status %s(%d): %v", want.String(), want, err)
			} else if got != want {
				t.Errorf("round-trip of %s(%d) came back as %s(%d)", want.String(), want, got.String(), got)
			}
		})
	}
}

// TestStatusOfString_Bad checks that trying to convert several bad statuses gives errors.
func TestStatusOfString_Bad(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"empty":   "",
		"unknown": "bad",
		"clipped": "timeou",
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := subject.StatusOfString(c)
			testhelp.ExpectErrorIs(t, err, subject.ErrBadStatus, "in bad StatusOfString")
		})
	}
}
