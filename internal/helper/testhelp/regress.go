// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package testhelp

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFilesOfExt is a scaffold for running a testing payload on every file in a directory dir that has extension ext.
// The tester function f receives the extensionless basename name and full path path.
// ext must contain a leading ".".
func TestFilesOfExt(t *testing.T, dir, ext string, f func(t *testing.T, name, path string)) {
	dfs := os.DirFS(filepath.ToSlash(dir))
	files, err := fs.Glob(dfs, "*"+ext)
	require.NoError(t, err, "couldn't stat test inputs")

	for _, file := range files {
		name := strings.TrimSuffix(file, ext)
		t.Run(name, func(t *testing.T) {
			f(t, name, filepath.Join(dir, file))
		})
	}
}
