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

// ExampleFromString is a runnable example for FromString.
func ExampleFromString() {
	fmt.Println(id.FromString("foo.bar.baz"))
	fmt.Println(id.FromString("FOO.BAR.BAZ"))
	fmt.Println(id.FromString("foo..bar.baz"))

	// Output:
	// foo.bar.baz
	// foo.bar.baz
	//
}

// ExampleID_IsEmpty is a runnable example for IsEmpty.
func ExampleID_IsEmpty() {
	fmt.Println(id.ID{}.IsEmpty())
	fmt.Println(id.FromString("").IsEmpty())
	fmt.Println(id.FromString("foo.bar.baz").IsEmpty())

	// Output:
	// true
	// true
	// false
}

// ExampleID_Join is a runnable example for Join.
func ExampleID_Join() {
	id1 := id.FromString("foo.bar")
	id2 := id.FromString("baz.barbaz")
	fmt.Println(id1.Join(id2).String())

	// empty IDs do nothing when joined
	fmt.Println(id.ID{}.Join(id1).String())
	fmt.Println(id2.Join(id.ID{}).String())

	// Output:
	// foo.bar.baz.barbaz
	// foo.bar
	// baz.barbaz
}

// ExampleID_Tags is a runnable example for Tags.
func ExampleID_Tags() {
	for _, tag := range id.FromString("foo.bar.baz").Tags() {
		fmt.Println(tag)
	}

	// Output:
	// foo
	// bar
	// baz
}

// ExampleID_Uncons is a runnable example for Uncons.
func ExampleID_Uncons() {
	_, _, ok := id.ID{}.Uncons()
	fmt.Println("uncons of empty ok?:", ok)

	// An uncons of a 1-tag ID returns that tag as the head.
	hd, tl, ok := id.FromString("foo").Uncons()
	fmt.Printf("foo: ok=%v, head=%q, tail=%q\n", ok, hd, tl)

	hd, tl, ok = id.FromString("foo.bar.baz").Uncons()
	fmt.Printf("foo.bar.baz: ok=%v, head=%q, tail=%q\n", ok, hd, tl)

	// Output:
	// uncons of empty ok?: false
	// foo: ok=true, head="foo", tail=""
	// foo.bar.baz: ok=true, head="foo", tail="bar.baz"
}

// ExampleID_Unsnoc is a runnable example for Unsnoc.
func ExampleID_Unsnoc() {
	_, _, ok := id.ID{}.Unsnoc()
	fmt.Println("unsnoc of empty ok?:", ok)

	// An unsnoc of a 1-tag ID returns that tag as the tail.
	hd, tl, ok := id.FromString("foo").Unsnoc()
	fmt.Printf("foo: ok=%v, head=%q, tail=%q\n", ok, hd, tl)

	hd, tl, ok = id.FromString("foo.bar.baz").Unsnoc()
	fmt.Printf("foo.bar.baz: ok=%v, head=%q, tail=%q\n", ok, hd, tl)

	// Output:
	// unsnoc of empty ok?: false
	// foo: ok=true, head="", tail="foo"
	// foo.bar.baz: ok=true, head="foo.bar", tail="baz"
}

// ExampleID_Triple is a runnable example for Triple.
func ExampleID_Triple() {
	f, v, s := id.ID{}.Triple()
	fmt.Printf("empty ID: f=%q v=%q s=%q\n", f, v, s)

	f, v, s = id.FromString("x86").Triple()
	fmt.Printf("family ID: f=%q v=%q s=%q\n", f, v, s)

	f, v, s = id.FromString("x86.64").Triple()
	fmt.Printf("variant ID: f=%q v=%q s=%q\n", f, v, s)

	f, v, s = id.FromString("x86.64.coffeelake").Triple()
	fmt.Printf("subvariant ID: f=%q v=%q s=%q\n", f, v, s)

	// Output:
	// empty ID: f="" v="" s=""
	// family ID: f="x86" v="" s=""
	// variant ID: f="x86" v="64" s=""
	// subvariant ID: f="x86" v="64" s="coffeelake"
}

// TestNew_valid tests New using various 'valid' inputs.
func TestNew_valid(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		tags []string
		want string
	}{
		"empty":      {tags: []string{""}},
		"one-tag":    {tags: []string{"foo"}, want: "foo"},
		"multi-tag":  {tags: []string{"foo", "bar", "baz"}, want: "foo.bar.baz"},
		"hyphenated": {tags: []string{"weird-hyphens", "allowed"}, want: "weird-hyphens.allowed"},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if d, err := id.New(c.tags...); err != nil {
				t.Errorf("New from tags %v error: %v", c.tags, err)
			} else if d.String() != c.want {
				t.Errorf("New from tags %v=%s, want %s", c.tags, d.String(), c.want)
			}
		})
	}
}

// TestNew_valid tests New using various 'erroneous' inputs.
func TestNew_errors(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		tags []string
		want error
	}{
		"empty": {tags: []string{"foo", "", "bar"}, want: id.ErrTagEmpty},
		"sep":   {tags: []string{"oh.no", "spaghetti.o"}, want: id.ErrTagHasSep},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := id.New(c.tags...)
			testhelp.ExpectErrorIs(t, err, c.want, "New on erroneous tags")
		})
	}
}
