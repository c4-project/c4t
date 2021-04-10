// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package filekind_test

import (
	"fmt"
	"testing"

	"github.com/c4-project/c4t/internal/helper/testhelp"

	"github.com/c4-project/c4t/internal/model/filekind"
	"github.com/stretchr/testify/assert"
)

// ExampleKind_Strings is a runnable example for Strings.
func ExampleKind_Strings() {
	for _, s := range (filekind.C | filekind.Litmus | filekind.Trace).Strings() {
		fmt.Println(s)
	}

	// Unordered output:
	// c
	// litmus
	// trace
}

// ExampleKindFromStrings is a runnable example for KindFromStrings.
func ExampleKindFromStrings() {
	k, err := filekind.KindFromStrings("other", "c")
	if err != nil {
		fmt.Println("ERROR:", err)
	} else {
		fmt.Println(k.String())
	}

	// Output:
	// other|c
}

// TestKind_Matches tests various combinations of kind matching.
func TestKind_Matches(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		matchee, matcher filekind.Kind
		want             bool
	}{
		"bin/any": {
			matchee: filekind.Bin,
			matcher: filekind.Any,
			want:    true,
		},
		"trace/trace": {
			matchee: filekind.Trace,
			matcher: filekind.Trace,
			want:    true,
		},
		"litmus/csrc": {
			matchee: filekind.Litmus,
			matcher: filekind.CSrc,
			want:    false,
		},
		"csrc/c": {
			matchee: filekind.CSrc,
			matcher: filekind.C,
			want:    true,
		},
		"cheader/c": {
			matchee: filekind.CHeader,
			matcher: filekind.C,
			want:    true,
		},
		"c/csrc": {
			matchee: filekind.C,
			matcher: filekind.CSrc,
			want:    false,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := c.matchee.Matches(c.matcher)
			assert.Equalf(t, c.want, got, "%d matching %d", c.matchee, c.matcher)
		})
	}
}

// TestKind_MarshalJSON_roundTrip tests the JSON marshalling of Kind by round-tripping.
func TestKind_MarshalJSON_roundTrip(t *testing.T) {
	t.Parallel()

	cases := []filekind.Kind{
		0,
		filekind.Trace,
		filekind.CHeader,
		filekind.C,
		filekind.Bin | filekind.Litmus,
		filekind.C | filekind.Log,
	}
	for _, c := range cases {
		c := c
		t.Run(c.String(), func(t *testing.T) {
			t.Parallel()
			testhelp.TestJSONRoundTrip(t, c, "filekind")
		})
	}
}
