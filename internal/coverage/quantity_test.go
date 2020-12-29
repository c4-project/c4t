// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package coverage_test

import (
	"fmt"
	"testing"

	"github.com/c4-project/c4t/internal/coverage"
	"github.com/stretchr/testify/assert"
)

// ExampleQuantitySet_Buckets is a runnable example for QuantitySet.Buckets.
func ExampleQuantitySet_Buckets() {
	qs := coverage.QuantitySet{Count: 1000, Divisions: []int{4, 5}}
	for _, b := range qs.Buckets() {
		fmt.Println(b)
	}

	// Output:
	// 4[250]
	// 3[250]
	// 2[250]
	// 1_5[50]
	// 1_4[50]
	// 1_3[50]
	// 1_2[50]
	// 1_1[50]

}

// ExampleQuantitySet_Buckets is a runnable example for QuantitySet.Buckets where there is no division.
func ExampleQuantitySet_Buckets_noDivision() {
	qs := coverage.QuantitySet{Count: 1000, Divisions: []int{}}
	for _, b := range qs.Buckets() {
		fmt.Println(b)
	}

	// Output:
	// 1[1000]
}

// ExampleQuantitySet_Buckets_uneven is a runnable example for QuantitySet.Buckets when there is uneven division.
func ExampleQuantitySet_Buckets_uneven() {
	qs := coverage.QuantitySet{Count: 1000, Divisions: []int{3, 3}}
	for _, b := range qs.Buckets() {
		fmt.Println(b)
	}

	// Output:
	// 3[333]
	// 2[333]
	// 1_3[111]
	// 1_2[111]
	// 1_1[112]
}

func TestQuantitySet_Override(t *testing.T) {
	t.Parallel()

	base := coverage.QuantitySet{Count: 6, Divisions: []int{2, 4, 6, 8}}

	cases := map[string]struct {
		override, want coverage.QuantitySet
	}{
		"empty": {
			override: coverage.QuantitySet{},
			want:     base,
		},
		"nil-div": {
			override: coverage.QuantitySet{Count: 42, Divisions: nil},
			want:     coverage.QuantitySet{Count: 42, Divisions: []int{2, 4, 6, 8}},
		},
		"empty-div": {
			override: coverage.QuantitySet{Count: 42, Divisions: []int{}},
			want:     coverage.QuantitySet{Count: 42, Divisions: []int{2, 4, 6, 8}},
		},
		"zero-count": {
			override: coverage.QuantitySet{Count: 0, Divisions: []int{3, 5, 7, 9}},
			want:     coverage.QuantitySet{Count: 6, Divisions: []int{3, 5, 7, 9}},
		},
		"all": {
			override: coverage.QuantitySet{Count: 42, Divisions: []int{3, 5, 7, 9}},
			want:     coverage.QuantitySet{Count: 42, Divisions: []int{3, 5, 7, 9}},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			got := base
			got.Override(c.override)
			assert.Equal(t, c.want, got, "unexpected override result")
		})
	}
}

// TestQuantitySet_Buckets tests various corner cases of QuantitySet.Buckets.
func TestQuantitySet_Buckets(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		n   int
		in  []int
		out []coverage.Bucket
	}{
		"nil":          {n: 1000, in: nil, out: []coverage.Bucket{{Name: "1", Size: 1000}}},
		"empty":        {n: 100, in: []int{}, out: []coverage.Bucket{{Name: "1", Size: 100}}},
		"one-div":      {n: 10, in: []int{1}, out: []coverage.Bucket{{Name: "1", Size: 10}}},
		"zero-divs":    {n: 2048, in: []int{0}, out: []coverage.Bucket{{Name: "1", Size: 2048}}},
		"neg-divs":     {n: 99, in: []int{-2}, out: []coverage.Bucket{{Name: "1", Size: 99}}},
		"none+nil":     {n: 0, in: nil, out: []coverage.Bucket{}},
		"none+empty":   {n: 0, in: []int{}, out: []coverage.Bucket{}},
		"none+one-div": {n: 0, in: []int{1}, out: []coverage.Bucket{}},
		"too-many":     {n: 1, in: []int{3}, out: []coverage.Bucket{{Name: "3", Size: 0}, {Name: "2", Size: 0}, {Name: "1", Size: 1}}},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			qs := coverage.QuantitySet{Count: c.n, Divisions: c.in}
			got := qs.Buckets()
			assert.ElementsMatch(t, c.out, got, "buckets not as expected")
		})
	}
}
