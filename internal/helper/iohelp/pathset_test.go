// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/c4-project/c4t/internal/helper/iohelp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMkdirs_rmdirs checks that iohelp.Mkdirs and iohelp.Rmdirs work in tandem to make and remove directories.
func TestMkdirs_rmdirs(t *testing.T) {
	root := t.TempDir()
	var dirs [10]string
	for i := range dirs {
		dirs[i] = filepath.Join(root, fmt.Sprintf("dir%d", i))
	}

	require.NoError(t, iohelp.Mkdirs(dirs[:]...), "couldn't make dirs")
	for _, d := range dirs {
		assert.DirExists(t, d, "dir not made", d)
	}
	require.NoError(t, iohelp.Rmdirs(dirs[:]...), "couldn't remove dirs")
	for _, d := range dirs {
		assert.NoDirExists(t, d, "dir not removed", d)
	}
}
