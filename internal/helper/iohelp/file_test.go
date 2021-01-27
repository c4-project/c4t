// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/c4-project/c4t/internal/helper/iohelp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIsFileEmpty tests that IsFileEmpty works on various cases.
func TestIsFileEmpty(t *testing.T) {
	t.Parallel()
	td := t.TempDir()

	cases := map[string]struct {
		file string
		want bool
	}{
		"empty file":    {file: filepath.Join("testdata", "empty.txt"), want: true},
		"nonempty file": {file: filepath.Join("testdata", "nonempty.txt"), want: false},
		"new file":      {file: filepath.Join(td, "new.txt"), want: true},
	}

	for name, c := range cases {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f, err := os.OpenFile(c.file, os.O_CREATE|os.O_RDONLY, 0644)
			require.NoError(t, err, "opening/creating file shouldn't error")
			got, err := iohelp.IsFileEmpty(f)
			if assert.NoError(t, err, "IsFileEmpty shouldn't error") {
				assert.Equal(t, c.want, got, "IsFileEmpty returned wrong result")
			}
			require.NoError(t, f.Close(), "closing shouldn't error")
		})
	}
}
