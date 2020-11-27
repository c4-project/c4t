// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/MattWindsor91/c4t/internal/helper/testhelp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/c4t/internal/stage/fuzzer"
)

// TestParseSubjectCycle_roundTrip tests ParseSubjectCycle by round-tripping using String.
func TestParseSubjectCycle_roundTrip(t *testing.T) {
	t.Parallel()

	cases := map[string]fuzzer.SubjectCycle{
		"simple-zero":  {Name: "foo", Cycle: 0},
		"simple-one":   {Name: "bar", Cycle: 1},
		"complex-zero": {Name: "foo_bar", Cycle: 0},
		"complex-one":  {Name: "foo_bar", Cycle: 1},
	}

	for name, c := range cases {
		want := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := fuzzer.ParseSubjectCycle(want.String())

			require.NoError(t, err, "unexpected error in ParseSubjectCycle")
			assert.Equalf(t, want, got, "ParseSubjectCycle roundtrip diverge: got=%q, want=%q", got, want)
		})
	}
}

func TestParseSubjectCycle_errors(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  string
		err error
	}{
		"empty":    {in: "", err: fuzzer.ErrNotSubjectCycleName},
		"no-under": {in: "foobar", err: fuzzer.ErrNotSubjectCycleName},
		"no-num":   {in: "foo_bar", err: strconv.ErrSyntax},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := fuzzer.ParseSubjectCycle(c.in)
			testhelp.ExpectErrorIs(t, err, c.err, fmt.Sprintf("parsing bad subject-cycle %q", c.in))
		})
	}
}
