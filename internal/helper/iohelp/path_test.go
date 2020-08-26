// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

// ExampleExtlessFile is a runnable example for ExtlessFile.
func ExampleExtlessFile() {
	fmt.Println(iohelp.ExtlessFile("foo.c"))
	fmt.Println(iohelp.ExtlessFile("/home/piers/test"))
	fmt.Println(iohelp.ExtlessFile("/home/piers/example.txt"))

	// Output:
	// foo
	// test
	// example
}

// TestExpandMany_noExpand just tests that ExpandMany works properly when there is no need for expansion.
// (Testing when there _is_ is a bit unlikely to be robust.)
func TestExpandMany_noExpand(t *testing.T) {
	in := []string{"", "foo", filepath.Join("bar", "baz")}
	out, err := iohelp.ExpandMany(in)
	require.NoError(t, err, "expanding with no expansions shouldn't error")
	assert.ElementsMatch(t, in, out, "no expansion should have taken place")
}
