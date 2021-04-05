// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id_test

import (
	"fmt"
	"testing"

	"github.com/c4-project/c4t/internal/helper/testhelp"

	"github.com/c4-project/c4t/internal/id"
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
	c := map[id.ID]int{
		id.FromString("foo.bar"):       1,
		id.FromString("BAR"):           2,
		id.FromString("foobar.baz"):    3,
		id.FromString("barbaz.Foobaz"): 4,
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

// ExampleMapGlob is a runnable example for MapGlob.
func ExampleMapGlob() {
	c := map[id.ID]int{
		id.FromString("foo.baz"):     1,
		id.FromString("foo.bar.baz"): 2,
		id.FromString("foo.bar"):     3,
		id.FromString("bar.baz"):     4,
	}
	c2, _ := id.MapGlob(c, id.FromString("foo.*.baz"))
	for k, v := range c2.(map[id.ID]int) {
		fmt.Println(k, v)
	}

	// Unordered output:
	// foo.baz 1
	// foo.bar.baz 2
}

// ExampleLookupPrefix is a runnable example for LookupPrefix.
func ExampleLookupPrefix() {
	c := map[id.ID]int{
		id.FromString("foo"):     1,
		id.FromString("bar"):     2,
		id.FromString("bar.baz"): 3,
	}

	k1, v1, _ := id.LookupPrefix(c, id.FromString("bar.baz"))
	k2, v2, _ := id.LookupPrefix(c, id.FromString("bar.foobaz"))
	k3, v3, _ := id.LookupPrefix(c, id.FromString("foo.bar.baz"))

	fmt.Printf("matched bar.baz to %s (%d)\n", k1, v1)
	fmt.Printf("matched bar.foobaz to %s (%d)\n", k2, v2)
	fmt.Printf("matched foo.bar.baz to %s (%d)\n", k3, v3)

	// Output:
	// matched bar.baz to bar.baz (3)
	// matched bar.foobaz to bar (2)
	// matched foo.bar.baz to foo (1)
}

// ExampleSearchSlice is a runnable example for SearchSlice.
func ExampleSearchSlice() {
	haystack := []id.ID{
		id.FromString("fus"),
		id.FromString("fus.ro"),
		id.FromString("fus.ro.dah"),
	}

	fmt.Println(id.SearchSlice(haystack, id.ID{}))
	fmt.Println(id.SearchSlice(haystack, id.FromString("fus.dah")))
	fmt.Println(id.SearchSlice(haystack, id.FromString("fus.ro.dah")))
	fmt.Println(id.SearchSlice(haystack, id.FromString("fus.ro.dah.dah.dah")))

	// Output:
	// 0
	// 1
	// 2
	// 3
}

// TestMapGlob_errors tests MapKeys's error handling.
func TestMapGlob_errors(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in   interface{}
		glob id.ID
		out  error
	}{
		"normal": {
			in: map[id.ID]string{
				id.FromString("a"): "A",
				id.FromString("b"): "B",
				id.FromString("c"): "C",
			},
			glob: id.FromString("a.*"),
			out:  nil,
		},
		"not-a-map": {
			in:   5,
			glob: id.FromString("a.*"),
			out:  id.ErrNotMap,
		},
		"not-an-id-map": {
			in: map[string]int{
				"fus..ro": 6,
			},
			glob: id.FromString("a.*"),
			out:  id.ErrNotMap,
		},
		"glob-not-a-glob": {
			in: map[id.ID]int{
				id.FromString("a.b.c"): 6,
			},
			glob: id.FromString("a.*.b.*"),
			out:  id.ErrBadGlob,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := id.MapGlob(c.in, c.glob)
			testhelp.ExpectErrorIs(t, err, c.out, "testing MapGlob")
		})
	}
}

// TestMapKeys_errors tests MapKeys's error handling.
func TestMapKeys_errors(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  interface{}
		out error
	}{
		"normal": {
			in: map[id.ID]string{
				id.FromString("a"): "A",
				id.FromString("b"): "B",
				id.FromString("c"): "C",
			},
			out: nil,
		},
		"not-a-map": {
			in:  5,
			out: id.ErrNotMap,
		},
		"not-an-id-map": {
			in: map[string]int{
				"fus..ro": 6,
			},
			out: id.ErrNotMap,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := id.MapKeys(c.in)
			testhelp.ExpectErrorIs(t, err, c.out, "testing MapKeys")
		})
	}
}
