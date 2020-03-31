// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package gcc_test

import (
	"testing"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler/gcc"
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
