// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package id_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"
)

// ExampleSort is a runnable example for Sort.
func ExampleSort() {
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
	id.Sort(ids)
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

// ExampleMapKeys is a runnable example for MapKeys.
func ExampleMapKeys() {
	c := map[string]int{
		"foo.bar":       1,
		"BAR":           2,
		"foobar.baz":    3,
		"barbaz.Foobaz": 4,
	}
	ids, _ := id.MapKeys(c)
	for _, x := range ids {
		fmt.Println(x)
	}

	// Output:
	// bar
	// barbaz.foobaz
	// foo.bar
	// foobar.baz
}
