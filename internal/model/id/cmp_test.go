// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id_test

import (
	"fmt"
	"sort"

	"github.com/c4-project/c4t/internal/model/id"
)

// ExampleID_Less is a runnable example for Less.
func ExampleID_Less() {
	ids := []id.ID{
		id.FromString("arm.7"),
		id.FromString("arm.8"),
		id.FromString("ppc.64.le"),
		id.FromString("x86.32"),
		id.FromString("x86"),
		id.FromString("arm"),
		id.FromString("ppc"),
		id.FromString("x86.64"),
		id.FromString("ppc.64"),
		id.FromString("arm.6"),
	}
	// Note: in general, use id.Sort instead!
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].Less(ids[j])
	})
	for _, i := range ids {
		fmt.Println(i)
	}

	// Output:
	// arm
	// arm.6
	// arm.7
	// arm.8
	// ppc
	// ppc.64
	// ppc.64.le
	// x86
	// x86.32
	// x86.64
}

// ExampleID_Less is a runnable example for Equal.
func ExampleID_Equal() {
	fmt.Println(id.FromString("arm.7").Equal(id.FromString("arm.7")))
	fmt.Println(id.FromString("arm.7").Equal(id.FromString("arm.8")))
	fmt.Println(id.FromString("arm.7").Equal(id.FromString("ARM.8")))
	fmt.Println(id.FromString("arm.7").Equal(id.FromString("arm")))
	fmt.Println(id.ID{}.Equal(id.FromString("")))

	// Output:
	// true
	// false
	// false
	// false
	// true
}
