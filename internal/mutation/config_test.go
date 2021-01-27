// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/mutation"
)

// ExampleAutoConfig_Mutants is a runnable example for Config.Mutants.
func ExampleAutoConfig_Mutants() {
	cfg := mutation.AutoConfig{
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

// ExampleAutoConfig_IsActive is a runnable example for Config.IsActive.
func ExampleAutoConfig_IsActive() {
	cfg := mutation.AutoConfig{
		Ranges: []mutation.Range{{Start: 1, End: 4}},
	}

	fmt.Println("disabled with ranges:", cfg.IsActive())

	cfg.ChangeKilled = true
	fmt.Println("after-killed with ranges:", cfg.IsActive())

	cfg.ChangeKilled = false
	cfg.ChangeAfter = quantity.Timeout(1 * time.Minute)
	fmt.Println("after-time with ranges:", cfg.IsActive())

	cfg.Ranges[0].Start = 4
	fmt.Println("after-time with bad ranges:", cfg.IsActive())

	// Output:
	// disabled with ranges: false
	// after-killed with ranges: true
	// after-time with ranges: true
	// after-time with bad ranges: false
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

// ExampleRange_IsEmpty is a runnable example for IsEmpty.
func ExampleRange_IsEmpty() {
	fmt.Println("10..20:", mutation.Range{Start: 10, End: 20}.IsEmpty())
	fmt.Println("10..10:", mutation.Range{Start: 10, End: 10}.IsEmpty())
	fmt.Println("20..10:", mutation.Range{Start: 20, End: 10}.IsEmpty())
	fmt.Println("10..11:", mutation.Range{Start: 10, End: 11}.IsEmpty())

	// Output:
	// 10..20: false
	// 10..10: true
	// 20..10: true
	// 10..11: false
}

// TestRange_Mutants tests Range.Mutants with various cases.
func TestRange_Mutants(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  mutation.Range
		out []mutation.Mutant
	}{
		"empty": {
			in:  mutation.Range{},
			out: []mutation.Mutant{},
		},
		"singleton": {
			in:  mutation.Range{Start: 42, End: 43},
			out: []mutation.Mutant{42},
		},
		"inverted": {
			in:  mutation.Range{Start: 53, End: 27},
			out: []mutation.Mutant{},
		},
		"ok": {
			in:  mutation.Range{Start: 10, End: 20},
			out: []mutation.Mutant{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
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
