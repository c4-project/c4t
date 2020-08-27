// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/plan/analysis"
	"github.com/stretchr/testify/require"
)

// TestLoadFilterSet tests LoadFilterSet on an example filter file.
func TestLoadFilterSet(t *testing.T) {
	t.Parallel()

	got, err := analysis.LoadFilterSet(filepath.Join("testdata", "filters.yaml"))
	require.NoError(t, err, "loading filter set should not error")

	want := analysis.FilterSet{
		analysis.Filter{
			Style:             id.FromString("gcc.gcc"),
			MajorVersionBelow: 4,
			ErrorPattern:      "error: invalid memory model for ‘__atomic_exchange’",
		},
	}

	assert.ElementsMatch(t, got, want, "filter set not as expected")
}
