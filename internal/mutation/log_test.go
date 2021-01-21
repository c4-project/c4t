// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/mutation"
	"github.com/stretchr/testify/assert"
)

// TestScanLine tests ScanLine on various cases.
func TestScanLine(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		line     string
		hit, sel uint64
	}{
		"empty":       {line: "", hit: 0, sel: 0},
		"no-match":    {line: "the quick brown fox", hit: 0, sel: 0},
		"hit-missing": {line: "MUTATION HIT:", hit: 0, sel: 0},
		"hit-not-int": {line: "MUTATION HIT: Kappa", hit: 0, sel: 0},
		"hit-ok":      {line: "MUTATION HIT: 27", hit: 27, sel: 0},
		"hit-extra":   {line: "MUTATION HIT: 42 (not out)", hit: 42, sel: 0},
		"sel-missing": {line: "MUTATION SELECTED:", hit: 0, sel: 0},
		"sel-not-int": {line: "MUTATION SELECTED: Keepo", hit: 0, sel: 0},
		"sel-ok":      {line: "MUTATION SELECTED: 53", hit: 0, sel: 53},
		"sel-extra":   {line: "MUTATION SELECTED: 1990 (time for the Guru)", hit: 0, sel: 1990},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var hit, sel uint64

			mutation.ScanLine(c.line, func(u uint64) {
				hit = u
			},
				func(u uint64) {
					sel = u
				})

			assert.Equal(t, c.hit, hit, "hit mutant not updated")
			assert.Equal(t, c.sel, sel, "selected mutant not updated")
		})
	}
}
