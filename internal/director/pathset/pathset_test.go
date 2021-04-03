// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pathset_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/director/pathset"

	"github.com/c4-project/c4t/internal/id"
)

// ExampleNew is a runnable example for New.
func ExampleNew() {
	p := pathset.New(filepath.FromSlash("tests"))

	fmt.Println(filepath.ToSlash(p.DirScratch))
	fmt.Println(filepath.ToSlash(p.DirSaved))

	// Output:
	// tests/scratch
	// tests/saved
}

// ExamplePathset_Instance is a runnable example for Pathset.Instance.
func ExamplePathset_Instance() {
	p := pathset.Pathset{DirSaved: "saved", DirScratch: "scratch"}
	mid := id.FromString("foo.bar.baz")
	mi := p.Instance(mid)

	for _, path := range mi.Scratch.Dirs() {
		fmt.Println(filepath.ToSlash(path))
	}
	for _, path := range mi.Saved.DirList() {
		fmt.Println(filepath.ToSlash(path))
	}

	// Output:
	// scratch/foo/bar/baz/fuzz
	// scratch/foo/bar/baz/lift
	// scratch/foo/bar/baz/run
	// saved/foo/bar/baz/flagged
	// saved/foo/bar/baz/compile_fail
	// saved/foo/bar/baz/compile_timeout
	// saved/foo/bar/baz/run_fail
	// saved/foo/bar/baz/run_timeout
}

// TestPathset_Prepare tests Scratch.Prepare.
func TestPathset_Prepare(t *testing.T) {
	// Probably can't parallelise this - affects the filesystem?
	root := t.TempDir()
	p := pathset.New(root)

	assert.NoDirExists(t, p.DirSaved, "saved dir shouldn't exist yet")
	assert.NoDirExists(t, p.DirScratch, "saved dir shouldn't exist yet")
	require.NoError(t, p.Prepare(), "prepare shouldn't error on temp dir")
	assert.DirExists(t, p.DirSaved, "saved dir should now exist")
	assert.DirExists(t, p.DirScratch, "saved dir should now exist")
}
