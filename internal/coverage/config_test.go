// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/coverage"
	"github.com/stretchr/testify/require"
)

// TestLoadConfigFromFile loads a test coverage configuration from the test data, and compares it against expectations.
func TestLoadConfigFromFile(t *testing.T) {
	cfg, err := coverage.LoadConfigFromFile(filepath.Join("testdata", "coverage.toml"))
	require.NoError(t, err, "loading coverage file from testdata")

	assert.Equal(t, "~/coverage_out", cfg.Paths.OutDir, "outputs not as expected")
	assert.ElementsMatch(t, []string{"~/input"}, cfg.Paths.Inputs, "inputs not as expected")

	assert.Equal(t, 100_000, cfg.Quantities.Count, "count not as expected")
	assert.ElementsMatch(t, []int{10, 10}, cfg.Quantities.Divisions, "divisions not as expected")

	if assert.Contains(t, cfg.Profiles, "csmith", "profiles should contain csmith") {
		p := cfg.Profiles["csmith"]
		assert.Equal(t, coverage.Standalone, p.Kind, "csmith profile kind not as expected")
		if assert.NotNil(t, p.Run, "csmith run info not present") {
			assert.Equal(t, "csmith", p.Run.Cmd, "csmith command not as expected")
			assert.ElementsMatch(t, []string{"-s", "${seed}", "-o", "${outputDir}/${i}.c", "${input}"}, p.Run.Args, "csmith args not as expected")
		}
	}
}
