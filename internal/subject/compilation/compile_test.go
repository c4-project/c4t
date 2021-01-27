// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compilation_test

import (
	"path"
	"testing"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/stretchr/testify/assert"
)

// TestCompileFileset_StripMissing tests that StripMissing works appropriately, with reference to a test filesystem.
func TestCompileFileset_StripMissing(t *testing.T) {
	// See testdata/strip_missing.
	cf := compilation.CompileFileset{
		Bin: path.Join("testdata", "strip_missing", "a.out"),   // missing
		Log: path.Join("testdata", "strip_missing", "log.txt"), // not missing
	}

	want := cf
	want.Bin = ""

	got := cf.StripMissing()
	assert.Equal(t, got, want, "StripMissing (bin should be missing)")
}
