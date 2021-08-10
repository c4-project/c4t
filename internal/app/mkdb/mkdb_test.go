// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mkdb_test

import (
	"io"
	"path/filepath"
	"testing"

	"github.com/c4-project/c4t/internal/app/mkdb"
	"github.com/stretchr/testify/assert"
)

// TestApp tests that sql construction works on a scratch file.
func TestApp_printGlobalPath(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "c4t.db")
	args := []string{
		mkdb.Name,
		"-" + mkdb.FlagDbPath,
		path,
	}
	assert.NoError(t, mkdb.App(io.Discard, io.Discard).Run(args))
	// TODO(@MattWindsor91): check database is set-up.
}
