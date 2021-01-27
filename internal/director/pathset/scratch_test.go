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
)

// ExampleNewScratch is a runnable example for NewScratch.
func ExampleNewScratch() {
	p := pathset.NewScratch("scratch")

	fmt.Println("run: ", filepath.ToSlash(p.DirRun))
	fmt.Println("lift:", filepath.ToSlash(p.DirLift))
	fmt.Println("fuzz:", filepath.ToSlash(p.DirFuzz))

	// Output:
	// run:  scratch/run
	// lift: scratch/lift
	// fuzz: scratch/fuzz
}

// TestScratch_Prepare tests Scratch.Prepare.
func TestScratch_Prepare(t *testing.T) {
	// Probably can't parallelise this - affects the filesystem?
	root := t.TempDir()
	p := pathset.NewScratch(root)

	for _, d := range p.Dirs() {
		assert.NoDirExists(t, d, "dir shouldn't exist yet")
	}
	require.NoError(t, p.Prepare(), "prepare shouldn't error on temp dir")
	for _, d := range p.Dirs() {
		assert.DirExists(t, d, "dir should now exist")
	}
}
