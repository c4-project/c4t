// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package status_test

import (
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/stretchr/testify/assert"
)

// TestFlag_Matches tests the behaviour of Matches on various input pairs.
func TestFlag_Matches(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		toMatch, matchAgainst status.Flag
		want                  bool
	}{
		"ok against ok": {
			toMatch:      status.FlagOk,
			matchAgainst: status.FlagOk,
			want:         true,
		},
		"ok against flagged": {
			toMatch:      status.FlagOk,
			matchAgainst: status.FlagFlagged,
			want:         false,
		},
		"flagged against flagged": {
			toMatch:      status.FlagFlagged,
			matchAgainst: status.FlagFlagged,
			want:         true,
		},
		"fails against c-fail": {
			toMatch:      status.FlagFail,
			matchAgainst: status.FlagCompileFail,
			want:         true,
		},
		"c-fail against fails": {
			toMatch:      status.FlagCompileFail,
			matchAgainst: status.FlagFail,
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
