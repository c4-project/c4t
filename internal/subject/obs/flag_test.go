// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package obs_test

import (
	"fmt"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/subject/obs"
)

// ExampleFlagOfStrings is a testable example for ObsFlagOfStrings.
func ExampleFlagOfStrings() {
	f, _ := obs.FlagOfStrings("unsat", "undef")
	fmt.Println(f.Has(obs.Sat))
	fmt.Println(f.Has(obs.Unsat))
	fmt.Println(f.Has(obs.Undef))

	// Output:
	// false
	// true
	// true
}

// ExampleFlag_Strings is a testable example for Flag.Strings.
func ExampleFlag_Strings() {
	for _, s := range (obs.Sat | obs.Undef).Strings() {
		fmt.Println(s)
	}

	// Output:
	// sat
	// undef
}

// TestFlagOfStrings tests various cases of FlagOfStrings.
func TestFlagOfStrings(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  []string
		out obs.Flag
		err error
	}{
		"sat": {
			in:  []string{"sat"},
			out: obs.Sat,
		},
		"unsat": {
			in:  []string{"unsat"},
			out: obs.Unsat,
		},
		"undef": {
			in:  []string{"undef"},
			out: obs.Undef,
		},
		"sat-undef": {
			in:  []string{"sat", "undef"},
			out: obs.Sat | obs.Undef,
		},
		"unsat-undef": {
			in:  []string{"unsat", "undef"},
			out: obs.Unsat | obs.Undef,
		},
		"unknown": {
			in:  []string{"blurble"},
			err: obs.ErrBadFlag,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out, err := obs.FlagOfStrings(c.in...)
			if testhelp.ExpectErrorIs(t, err, c.err, "FlagOfStrings") {
				if err == nil {
					assert.Equal(t, c.out, out, "FlagOfStrings on:", c.in)
				}
			}
		})
	}
}
