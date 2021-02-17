// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/helper/iohelp"
)

// TestEnsureWriter tests various properties of EnsureWriter.
func TestEnsureWriter(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, io.Discard, iohelp.EnsureWriter(nil), "EnsureWriter(nil) didn't return discard")
	})
	t.Run("non-nil", func(t *testing.T) {
		t.Parallel()
		var b bytes.Buffer
		require.Equal(t, &b, iohelp.EnsureWriter(&b), "EnsureWriter(b) changed pointer")
	})
}
