// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stringhelp_test

import (
	"fmt"
	"testing"

	"github.com/c4-project/c4t/internal/helper/stringhelp"
)

// ExamplePluralQuantity is a runnable example for PluralQuantity.
func ExamplePluralQuantity() {
	fmt.Println(stringhelp.PluralQuantity(0, "director", "y", "ies"))
	fmt.Println(stringhelp.PluralQuantity(1, "fil", "e", "es"))
	fmt.Println(stringhelp.PluralQuantity(2, "corp", "us", "ora"))

	// Output:
	// 0 directories
	// 1 file
	// 2 corpora
}

// BenchmarkPluralQuantity_zero benchmarks PluralQuantity with a quantity of 0.
func BenchmarkPluralQuantity_zero(b *testing.B) {
	benchmarkPluralQuantity(0, b)
}

// BenchmarkPluralQuantity_one benchmarks PluralQuantity with a quantity of 1.
func BenchmarkPluralQuantity_one(b *testing.B) {
	benchmarkPluralQuantity(1, b)
}

// BenchmarkPluralQuantity_more benchmarks PluralQuantity with a quantity of more than 1.
func BenchmarkPluralQuantity_more(b *testing.B) {
	benchmarkPluralQuantity(2, b)
}

func benchmarkPluralQuantity(n int, b *testing.B) {
	var result string

	for i := 0; i < b.N; i++ {
		result = stringhelp.PluralQuantity(n, "director", "y", "ies")
	}

	// This just exists to stop 'result' being flagged as unused.
	if result == "" {
		b.Fail()
	}
}
