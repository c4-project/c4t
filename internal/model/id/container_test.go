// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id_test

import (
	"fmt"
	"testing"

	"github.com/c4-project/c4t/internal/helper/stringhelp"

	"github.com/c4-project/c4t/internal/helper/testhelp"

	"github.com/c4-project/c4t/internal/model/id"
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

// ExampleMapGlob is a runnable example for MapGlob.
func ExampleMapGlob() {
	c := map[string]int{
		"foo.baz":     1,
		"foo.bar.baz": 2,
		"foo.bar":     3,
		"bar.baz":     4,
	}
	c2, _ := id.MapGlob(c, id.FromString("foo.*.baz"))
	for k, v := range c2.(map[string]int) {
		fmt.Println(k, v)
	}

	// Unordered output:
	// foo.baz 1
	// foo.bar.baz 2
}

func ExampleLookupPrefix() {
	c := map[string]int{
		"foo":     1,
		"bar":     2,
		"bar.baz": 3,
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

// TestMapGlob_errors tests MapKeys's error handling.
func TestMapGlob_errors(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in   interface{}
		glob id.ID
		out  error
	}{
		"normal": {
			in: map[string]string{
				"a": "A",
				"b": "B",
				"c": "C",
			},
			glob: id.FromString("a.*"),
			out:  nil,
		},
		"not-a-map": {
			in:   5,
			glob: id.FromString("a.*"),
			out:  stringhelp.ErrNotMap,
		},
		"key-not-an-id": {
			in: map[string]int{
				"fus..ro": 6,
			},
			glob: id.FromString("a.*"),
			out:  id.ErrTagEmpty,
		},
		"glob-not-a-glob": {
			in: map[string]int{
				"a.b.c": 6,
			},
			glob: id.FromString("a.*.b.*"),
			out:  id.ErrBadGlob,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
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
			in: map[string]string{
				"a": "A",
				"b": "B",
				"c": "C",
			},
			out: nil,
		},
		"not-a-map": {
			in:  5,
			out: stringhelp.ErrNotMap,
		},
		"not-an-id": {
			in: map[string]int{
				"fus..ro": 6,
			},
			out: id.ErrTagEmpty,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			_, err := id.MapKeys(c.in)
			testhelp.ExpectErrorIs(t, err, c.out, "testing MapKeys")
		})
	}
}
