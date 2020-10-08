// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import "fmt"

// ExampleQuantitySet_Buckets is a runnable example for QuantitySet.Buckets.
func ExampleQuantitySet_Buckets() {
	qs := QuantitySet{Count: 1000, Divisions: []int{4, 5}}
	for bname, bsize := range qs.Buckets() {
		fmt.Printf("%q: %d\n", bname, bsize)
	}

	// Unordered output:
	// "1,1": 50
	// "1,2": 50
	// "1,3": 50
	// "1,4": 50
	// "1,5": 50
	// "2": 250
	// "3": 250
	// "4": 250
}

// ExampleQuantitySet_Buckets is a runnable example for QuantitySet.Buckets where there is no division.
func ExampleQuantitySet_Buckets_noDivision() {
	qs := QuantitySet{Count: 1000, Divisions: []int{}}
	for bname, bsize := range qs.Buckets() {
		fmt.Printf("%q: %d\n", bname, bsize)
	}

	// Output:
	// "1": 1000
}

// ExampleQuantitySet_Buckets_uneven is a runnable example for QuantitySet.Buckets when there is uneven division.
func ExampleQuantitySet_Buckets_uneven() {
	qs := QuantitySet{Count: 1000, Divisions: []int{3, 3}}
	for bname, bsize := range qs.Buckets() {
		fmt.Printf("%q: %d\n", bname, bsize)
	}

	// Unordered output:
	// "1,1": 112
	// "1,2": 111
	// "1,3": 111
	// "2": 333
	// "3": 333
}
