// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/mutation"
)

// ExampleConfig_Mutants is a runnable example for Config.
func ExampleConfig_Mutants() {
	cfg := mutation.Config{
		Ranges: []mutation.Range{
			{Start: 1, End: 4},
			{Start: 10, End: 11},
			{Start: 27, End: 31},
		},
	}

	fmt.Print("mutants:")
	for _, i := range cfg.Mutants() {
		fmt.Printf(" %d", i)
	}
	fmt.Println()

	// Output:
	// mutants: 1 2 3 10 27 28 29 30
}

// ExampleRange_Mutants is a runnable example for Range.
func ExampleRange_Mutants() {
	fmt.Print("10..20:")
	for _, i := range (mutation.Range{Start: 10, End: 20}).Mutants() {
		fmt.Printf(" %d", i)
	}
	fmt.Println()

	// Output:
	// 10..20: 10 11 12 13 14 15 16 17 18 19
}

// TestRange_Mutants tests Range.Mutants with various cases.
func TestRange_Mutants(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  mutation.Range
		out []uint
	}{
		"empty": {
			in:  mutation.Range{},
			out: []uint{},
		},
		"singleton": {
			in:  mutation.Range{Start: 42, End: 43},
			out: []uint{42},
		},
		"inverted": {
			in:  mutation.Range{Start: 53, End: 27},
			out: []uint{},
		},
		"ok": {
			in:  mutation.Range{Start: 10, End: 20},
			out: []uint{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, c.out, c.in.Mutants())
		})
	}
}
