// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package status_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/stretchr/testify/assert"
)

// TestFlag_MatchesAll tests the behaviour of Flag.MatchesAll on various input pairs.
func TestFlag_MatchesAll(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		toMatch, matchAgainst status.Flag
		want                  bool
	}{
		"flagged against flagged": {
			toMatch:      status.FlagFlagged,
			matchAgainst: status.FlagFlagged,
			want:         true,
		},
		"flagged against c-fail": {
			toMatch:      status.FlagFlagged,
			matchAgainst: status.FlagCompileFail,
			want:         false,
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

			assert.Equal(t, c.want, c.toMatch.MatchesAll(c.matchAgainst))
		})
	}
}

// TestFlag_MatchesAny tests the behaviour of Flag.MatchesAny on various input pairs.
func TestFlag_MatchesAny(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		toMatch, matchAgainst status.Flag
		want                  bool
	}{
		"flagged against flagged": {
			toMatch:      status.FlagFlagged,
			matchAgainst: status.FlagFlagged,
			want:         true,
		},
		"flagged against c-fail": {
			toMatch:      status.FlagFlagged,
			matchAgainst: status.FlagCompileFail,
			want:         false,
		},
		"fails against c-fail": {
			toMatch:      status.FlagFail,
			matchAgainst: status.FlagCompileFail,
			want:         true,
		},
		"c-fail against fails": {
			toMatch:      status.FlagCompileFail,
			matchAgainst: status.FlagFail,
			want:         true,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, c.want, c.toMatch.MatchesAny(c.matchAgainst))
		})
	}
}

// TestFlag_MatchesStatus tests the behaviour of Flag.MatchesStatus on various input pairs.
func TestFlag_MatchesStatus(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		toMatch      status.Flag
		matchAgainst status.Status
		want         bool
	}{
		"flagged against ok": {
			toMatch:      status.FlagFlagged,
			matchAgainst: status.Ok,
			want:         false,
		},
		"filtered against ok": {
			toMatch:      status.FlagFiltered,
			matchAgainst: status.Ok,
			want:         false,
		},
		"flagged against flagged": {
			toMatch:      status.FlagFlagged,
			matchAgainst: status.Flagged,
			want:         true,
		},
		"fails against c-fail": {
			toMatch:      status.FlagFail,
			matchAgainst: status.CompileFail,
			want:         true,
		},
		"filtered fails against c-fail": {
			toMatch:      status.FlagFail | status.FlagFiltered,
			matchAgainst: status.CompileFail,
			want:         true,
		},
		"c-fails against filtered": {
			toMatch:      status.FlagCompileFail,
			matchAgainst: status.Filtered,
			want:         false,
		},
		"filtered c-fails against filtered": {
			toMatch:      status.FlagCompileFail | status.FlagFiltered,
			matchAgainst: status.Filtered,
			want:         true,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, c.want, c.toMatch.MatchesStatus(c.matchAgainst))
		})
	}
}
