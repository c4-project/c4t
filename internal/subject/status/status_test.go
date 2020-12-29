// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package status_test

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/helper/testhelp"
)

// ExampleStatus_IsOk is a runnable example for Status.IsOk.
func ExampleStatus_IsOk() {
	fmt.Println("is", status.Ok, "ok?", status.Ok.IsOk())
	fmt.Println("is", status.Flagged, "ok?", status.Flagged.IsOk())
	fmt.Println("is", status.CompileFail, "ok?", status.CompileFail.IsOk())

	// Output:
	// is Ok ok? true
	// is Flagged ok? false
	// is CompileFail ok? false
}

// TestFromCompileError tests several run errors to see if their status equivalent is as expected.
func TestFromCompileError(t *testing.T) {
	t.Parallel()

	e := errors.New("bloop")

	cases := map[string]struct {
		in   error
		want status.Status
		out  error
	}{
		"ok": {
			in:   nil,
			want: status.Ok,
		},
		"timeout": {
			in:   context.DeadlineExceeded,
			want: status.CompileTimeout,
		},
		"err": {
			in:   &exec.ExitError{},
			want: status.CompileFail,
		},
		"other": {
			in:  e,
			out: e,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := status.FromCompileError(c.in)
			testhelp.ExpectErrorIs(t, err, c.out, "FromCompileError foldback")
			assert.Equal(t, c.want, got, "unexpected status")
		})
	}
}

// TestFromRunError tests several run errors to see if their status equivalent is as expected.
func TestFromRunError(t *testing.T) {
	t.Parallel()

	e := errors.New("bloop")

	cases := map[string]struct {
		in   error
		want status.Status
		out  error
	}{
		"ok": {
			in:   nil,
			want: status.Ok,
		},
		"timeout": {
			in:   context.DeadlineExceeded,
			want: status.RunTimeout,
		},
		"err": {
			in:   &exec.ExitError{},
			want: status.RunFail,
		},
		"other": {
			in:  e,
			out: e,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := status.FromRunError(c.in)
			testhelp.ExpectErrorIs(t, err, c.out, "FromRunError foldback")
			assert.Equal(t, c.want, got, "unexpected status")
		})
	}
}

// TestFromString_roundTrip checks that converting a status to and back from its string is the identity.
func TestFromString_roundTrip(t *testing.T) {
	t.Parallel()
	for want := status.Unknown; want <= status.Last; want++ {
		want := want
		t.Run(strconv.Itoa(int(want)), func(t *testing.T) {
			t.Parallel()
			got, err := status.FromString(want.String())
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
			_, err := status.FromString(c)
			testhelp.ExpectErrorIs(t, err, status.ErrBad, "in bad FromString")
		})
	}
}
