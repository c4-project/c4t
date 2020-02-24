package runner_test

import (
	"strconv"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/testhelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/runner"
)

// TestStatusOfString_RoundTrip checks that converting a status to and back from its string is the identity.
func TestStatusOfString_RoundTrip(t *testing.T) {
	for want := runner.StatusUnknown; want < runner.NumStatus; want++ {
		t.Run(strconv.Itoa(int(want)), func(t *testing.T) {
			got, err := runner.StatusOfString(want.String())
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
	cases := map[string]string{
		"empty":   "",
		"unknown": "bad",
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := runner.StatusOfString(c)
			testhelp.ExpectErrorIs(t, err, runner.ErrBadStatus, "in bad StatusOfString")
		})
	}
}
