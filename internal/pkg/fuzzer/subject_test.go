// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer_test

import (
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"
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
			got, err := fuzzer.ParseSubjectCycle(want.String())
			if err != nil {
				t.Fatal("unexpected error in ParseSubjectCycle:", err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("ParseSubjectCycle roundtrip diverge: got=%q, want=%q", got, want)
			}
		})
	}
}
