// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id_test

import (
	"fmt"
	"testing"

	"github.com/MattWindsor91/c4t/internal/helper/testhelp"
	"github.com/MattWindsor91/c4t/internal/model/id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ExampleID_HasPrefix is a runnable example for HasPrefix.
func ExampleID_HasPrefix() {
	x := id.FromString("x86.64")
	fmt.Println("x86.64 prefix of x86.64:", x.HasPrefix(id.FromString("x86.64")))
	fmt.Println("x86 prefix of x86.64:", x.HasPrefix(id.FromString("x86")))
	fmt.Println("arm prefix of x86.64:", x.HasPrefix(id.FromString("arm")))
	fmt.Println("empty prefix of x86.64:", x.HasPrefix(id.ID{}))

	// Output:
	// x86.64 prefix of x86.64: true
	// x86 prefix of x86.64: true
	// arm prefix of x86.64: false
	// empty prefix of x86.64: true
}

// ExampleID_HasSuffix is a runnable example for HasSuffix.
func ExampleID_HasSuffix() {
	x := id.FromString("x86.64")
	fmt.Println("x86.64 suffix of x86.64:", x.HasSuffix(id.FromString("x86.64")))
	fmt.Println("64 suffix of x86.64:", x.HasSuffix(id.FromString("64")))
	fmt.Println("32 suffix of x86.64:", x.HasSuffix(id.FromString("32")))
	fmt.Println("empty suffix of x86.64:", x.HasSuffix(id.ID{}))

	// Output:
	// x86.64 suffix of x86.64: true
	// 64 suffix of x86.64: true
	// 32 suffix of x86.64: false
	// empty suffix of x86.64: true
}

// TestID_HasPrefix tests the HasPrefix function through various cases.
func TestID_HasPrefix(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inTags     []string
		prefixTags []string
		want       bool
	}{
		"yes-empty-empty": {
			inTags:     []string{},
			prefixTags: []string{},
			want:       true,
		},
		"yes-nonempty-empty": {
			inTags:     []string{"arm", "8", "a"},
			prefixTags: []string{},
			want:       true,
		},
		"no-empty-nonempty": {
			inTags:     []string{},
			prefixTags: []string{"arm"},
			want:       false,
		},
		"yes-prefix": {
			inTags:     []string{"arm", "8"},
			prefixTags: []string{"arm"},
			want:       true,
		},
		"no-prefix": {
			inTags:     []string{"arm", "8"},
			prefixTags: []string{"x86"},
			want:       false,
		},
		"yes-same-length": {
			inTags:     []string{"arm", "8"},
			prefixTags: []string{"arm", "8"},
			want:       true,
		},
		"no-same-length": {
			inTags:     []string{"arm", "8"},
			prefixTags: []string{"arm", "7"},
			want:       false,
		},
		"yes-same-length-single": {
			inTags:     []string{"arm"},
			prefixTags: []string{"arm"},
			want:       true,
		},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			in, err := id.New(c.inTags...)
			require.NoError(t, err)
			prefix, err := id.New(c.prefixTags...)
			require.NoError(t, err)

			got := in.HasPrefix(prefix)
			assert.Equal(t, c.want, got, "result of HasPrefix")
		})
	}
}

// TestID_HasSuffix tests the HasSuffix function through various cases.
func TestID_HasSuffix(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inTags     []string
		suffixTags []string
		want       bool
	}{
		"yes-empty-empty": {
			inTags:     []string{},
			suffixTags: []string{},
			want:       true,
		},
		"yes-nonempty-empty": {
			inTags:     []string{"arm", "8", "a"},
			suffixTags: []string{},
			want:       true,
		},
		"no-empty-nonempty": {
			inTags:     []string{},
			suffixTags: []string{"arm"},
			want:       false,
		},
		"yes-suffix": {
			inTags:     []string{"arm", "8"},
			suffixTags: []string{"8"},
			want:       true,
		},
		"no-suffix": {
			inTags:     []string{"arm", "8"},
			suffixTags: []string{"7"},
			want:       false,
		},
		"yes-same-length": {
			inTags:     []string{"arm", "8"},
			suffixTags: []string{"arm", "8"},
			want:       true,
		},
		"no-same-length": {
			inTags:     []string{"arm", "8"},
			suffixTags: []string{"arm", "7"},
			want:       false,
		},
		"yes-same-length-single": {
			inTags:     []string{"arm"},
			suffixTags: []string{"arm"},
			want:       true,
		},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			in, err := id.New(c.inTags...)
			require.NoError(t, err)
			suffix, err := id.New(c.suffixTags...)
			require.NoError(t, err)

			got := in.HasSuffix(suffix)
			assert.Equal(t, c.want, got, "result of HasSuffix")
		})
	}
}

// TestID_Matches tests the Matches function through various successful, failure, and error cases.
func TestID_Matches(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inTags   []string
		globTags []string
		want     bool
		err      error
	}{
		"yes-exact-empty": {
			inTags:   []string{},
			globTags: []string{},
			want:     true,
		},
		"yes-exact-nonempty": {
			inTags:   []string{"foo", "bar", "baz"},
			globTags: []string{"foo", "bar", "baz"},
			want:     true,
		},
		"no-exact-nonempty": {
			inTags:   []string{"foo", "bar", "baz"},
			globTags: []string{"foo", "baz", "bar"},
			want:     false,
		},
		"yes-glob-prefix": {
			inTags:   []string{"foo", "bar", "baz"},
			globTags: []string{"foo", id.TagGlob},
			want:     true,
		},
		"no-glob-prefix": {
			inTags:   []string{"foo", "bar", "baz"},
			globTags: []string{"baz", id.TagGlob},
			want:     false,
		},
		"yes-glob-suffix": {
			inTags:   []string{"foo", "bar", "baz"},
			globTags: []string{id.TagGlob, "bar", "baz"},
			want:     true,
		},
		"no-glob-suffix": {
			inTags:   []string{"foo", "bar", "baz"},
			globTags: []string{id.TagGlob, "bar"},
			want:     false,
		},
		"yes-glob-both-empty": {
			inTags:   []string{"alpha", "beta", "kappa", "lambda"},
			globTags: []string{"alpha", "beta", id.TagGlob, "kappa", "lambda"},
			want:     true,
		},
		"yes-glob-both-nonempty": {
			inTags:   []string{"alpha", "beta", "kappa", "lambda"},
			globTags: []string{"alpha", id.TagGlob, "lambda"},
			want:     true,
		},
		"no-glob-both": {
			inTags:   []string{"alpha", "beta", "kappa", "lambda"},
			globTags: []string{"alpha", "beta", id.TagGlob, "kappa"},
			want:     false,
		},
		"err-multiglob": {
			inTags:   []string{"alpha", "beta", "kappa", "lambda"},
			globTags: []string{"alpha", id.TagGlob, "beta", id.TagGlob, "lambda"},
			err:      id.ErrBadGlob,
		},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			in, err := id.New(c.inTags...)
			require.NoError(t, err)
			glob, err := id.New(c.globTags...)
			require.NoError(t, err)

			matches, err := in.Matches(glob)
			if testhelp.ExpectErrorIs(t, err, c.err, "Matches") {
				assert.Equal(t, c.want, matches, "result of Matches")
			}
		})
	}
}
