// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gcc_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/serviceimpl/compiler/gcc"
	"github.com/stretchr/testify/assert"
)

// TestOptLevelNames_consistency makes sure OptLevelNames is consistent with OptLevels in both directions.
func TestOptLevelNames_consistency(t *testing.T) {
	t.Parallel()
	t.Run("names-to-levels", func(t *testing.T) {
		t.Parallel()
		testNameConsistency(t, gcc.OptLevelNames)
	})
	t.Run("levels-to-names", func(t *testing.T) {
		t.Parallel()

		for n := range gcc.OptLevels {
			assert.Contains(t, gcc.OptLevelNames, n, "level not in names", n)
		}
	})
}

// TestGCC_Levels tests that OptLevels returns the expected level set.
func TestGCC_Levels(t *testing.T) {
	t.Parallel()
	ls, err := gcc.GCC{}.OptLevels(nil)
	assert.NoError(t, err)
	assert.Equal(t, gcc.OptLevels, ls)
}

// TestGCC_DefaultOptLevels tests that DefaultOptLevels returns a level set broadly consistent with expectations.
func TestGCC_DefaultOptLevels(t *testing.T) {
	t.Parallel()

	dl, err := gcc.GCC{}.DefaultOptLevels(nil)
	require.NoError(t, err)

	t.Run("disabled", func(t *testing.T) {
		t.Parallel()
		for _, d := range gcc.OptLevelDisabledNames {
			assert.NotContains(t, dl, d, "disabled opt level in defaults", d)
		}
	})
	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		for n := range dl {
			assert.Contains(t, gcc.OptLevels, n, "name not in levels", n)
		}
	})
}

// TestOptLevelDisabledNames_consistency makes sure OptLevelDisabledNames is consistent with OptLevels.
func TestOptLevelDisabledNames_consistency(t *testing.T) {
	t.Parallel()
	testNameConsistency(t, gcc.OptLevelDisabledNames)
}

// testNameConsistency tests that all of the names in names are in the OptLevels map
func testNameConsistency(t *testing.T, names []string) {
	t.Helper()
	for _, n := range names {
		assert.Contains(t, gcc.OptLevels, n, "name not in levels", n)
	}
}
