// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package status_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
)

// ExampleStatus_IsOk is a runnable example for IsOk.
func ExampleStatus_IsOk() {
	fmt.Println("is", status.Ok, "ok?", status.Ok.IsOk())
	fmt.Println("is", status.Flagged, "ok?", status.Flagged.IsOk())
	fmt.Println("is", status.CompileFail, "ok?", status.CompileFail.IsOk())

	// Output:
	// is ok ok? true
	// is flagged ok? false
	// is compile/fail ok? false
}

// TestOfString_RoundTrip checks that converting a status to and back from its string is the identity.
func TestOfString_RoundTrip(t *testing.T) {
	t.Parallel()
	for want := status.Unknown; want < status.Num; want++ {
		want := want
		t.Run(strconv.Itoa(int(want)), func(t *testing.T) {
			t.Parallel()
			got, err := status.OfString(want.String())
			if err != nil {
				t.Errorf("unexpected error round-tripping status %s(%d): %v", want.String(), want, err)
			} else if got != want {
				t.Errorf("round-trip of %s(%d) came back as %s(%d)", want.String(), want, got.String(), got)
			}
		})
	}
}

// TestOfString_Bad checks that trying to convert several bad statuses gives errors.
func TestOfString_Bad(t *testing.T) {
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
			_, err := status.OfString(c)
			testhelp.ExpectErrorIs(t, err, status.ErrBad, "in bad OfString")
		})
	}
}
