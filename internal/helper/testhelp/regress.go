// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package testhelp

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/1set/gut/yos"
	"github.com/stretchr/testify/require"
)

// TestFilesOfExt is a scaffold for running a testing payload on every file in a directory dir that has extension ext.
// The tester function f receives the extensionless basename name and full path path.
// ext must contain a leading ".".
func TestFilesOfExt(t *testing.T, dir, ext string, f func(t *testing.T, name, path string)) {
	ents, err := yos.ListMatch(dir, yos.ListIncludeFile, "*"+ext)
	require.NoError(t, err, "couldn't stat test inputs")

	for _, ent := range ents {
		path := ent.Path
		name := strings.TrimSuffix(filepath.Base(path), ext)
		t.Run(name, func(t *testing.T) {
			f(t, name, path)
		})
	}
}
