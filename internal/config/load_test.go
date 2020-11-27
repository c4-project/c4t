// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package config_test

import (
	"path/filepath"
	"testing"

	"github.com/MattWindsor91/c4t/internal/config"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

// TestLoad_direct tests config.Load with a directly provided filename.
func TestLoad_direct(t *testing.T) {
	conf, err := config.Load(filepath.Join("testdata", "tester.toml"))
	require.NoError(t, err, "error loading config file")

	assert.Equal(t, "/home/example/test_out", conf.Paths.OutDir, "OutDir not set correctly")
	assert.Equal(t, "/home/example/filters.yaml", conf.Paths.FilterFile, "FilterFile not set correctly")

	assert.NotNil(t, conf.Fuzz, "Fuzz not present")

	assert.ElementsMatch(t, []string{"/home/example/inputs", "/home/example/standalone.litmus"}, conf.Paths.Inputs, "Inputs not set correctly")
}

// TestLoad_broken tests config.Load with a bogus filename.
func TestLoad_broken(t *testing.T) {
	_, err := config.Load(filepath.Join("testdata", "not_toml.txt"))
	require.Error(t, err, "should error")
}
