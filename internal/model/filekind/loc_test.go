// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package filekind_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/c4t/internal/model/filekind"
)

// TestLoc_Matches tests various combinations of location matching.
func TestLoc_Matches(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		matchee, matcher filekind.Loc
		want             bool
	}{
		"fuzz/any": {
			matchee: filekind.InFuzz,
			matcher: filekind.Any,
			want:    true,
		},
		"recipe/recipe": {
			matchee: filekind.InRecipe,
			matcher: filekind.InRecipe,
			want:    true,
		},
		"orig/compile": {
			matchee: filekind.InOrig,
			matcher: filekind.InCompile,
			want:    false,
		},
		"fuzz/fuzz+orig": {
			matchee: filekind.InFuzz,
			matcher: filekind.InFuzz | filekind.InOrig,
			want:    true,
		},
		"fuzz+orig/fuzz": {
			matchee: filekind.InFuzz | filekind.InOrig,
			matcher: filekind.InFuzz,
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
