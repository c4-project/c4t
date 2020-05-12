// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"
)

// TestFlag_Matches tests the behaviour of Matches on various input pairs.
func TestFlag_Matches(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		toMatch, matchAgainst analysis.Flag
		want                  bool
	}{
		"ok against ok": {
			toMatch:      analysis.FlagOk,
			matchAgainst: analysis.FlagOk,
			want:         true,
		},
		"ok against flagged": {
			toMatch:      analysis.FlagOk,
			matchAgainst: analysis.FlagFlagged,
			want:         false,
		},
		"flagged against flagged": {
			toMatch:      analysis.FlagFlagged,
			matchAgainst: analysis.FlagFlagged,
			want:         true,
		},
		"fails against c-fail": {
			toMatch:      analysis.FlagFail,
			matchAgainst: analysis.FlagCompileFail,
			want:         true,
		},
		"c-fail against fails": {
			toMatch:      analysis.FlagCompileFail,
			matchAgainst: analysis.FlagFail,
			want:         false,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, c.want, c.toMatch.Matches(c.matchAgainst))
		})
	}
}
