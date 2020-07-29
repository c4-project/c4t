// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compilation_test

import (
	"path"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/subject/compilation"

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
