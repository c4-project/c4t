// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPathset_Prepare tests Pathset.Prepare in a temporary directory.
func TestPathset_Prepare(t *testing.T) {
	td := t.TempDir()
	ps := NewPathset(td)

	assert.NoDirExists(t, ps.DirLitmus)
	assert.NoDirExists(t, ps.DirTrace)

	err := ps.Prepare()
	require.NoError(t, err, "preparing fuzzer pathset in temp dir")

	assert.DirExists(t, ps.DirLitmus)
	assert.DirExists(t, ps.DirLitmus)
}
