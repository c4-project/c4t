// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package filekind_test

import (
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"
	"github.com/stretchr/testify/assert"
)

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
