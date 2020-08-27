// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis_test

import (
	"path/filepath"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/MattWindsor91/act-tester/internal/subject/status"

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

	want, werr := analysis.Compile(analysis.FilterSet{
		analysis.Filter{
			Style:             id.FromString("gcc"),
			MajorVersionBelow: 4,
			ErrorPattern:      "error: invalid memory model for ‘__atomic_exchange’",
		},
	})
	require.NoError(t, werr, "compiling filter set should not error")

	assert.ElementsMatch(t, got, want, "filter set not as expected")
}

func TestFilterSet_FilteredStatus(t *testing.T) {
	t.Parallel()

	fs, err := analysis.LoadFilterSet(filepath.Join("testdata", "filters.yaml"))
	require.NoError(t, err, "loading filter set should not error")

	cases := map[string]struct {
		inComp   compiler.Configuration
		inLog    string
		inResult compilation.CompileResult
		want     status.Status
	}{
		"no filtering": {
			inComp:   compiler.MockX86Gcc(),
			inLog:    "blep",
			inResult: compilation.CompileResult{Result: compilation.Result{Status: status.Ok}},
			want:     status.Ok,
		},
		"filtering in middle of string": {
			inComp:   compiler.MockX86Gcc(),
			inLog:    "foo error: invalid memory model for ‘__atomic_exchange’ bar",
			inResult: compilation.CompileResult{Result: compilation.Result{Status: status.Ok}},
			want:     status.Filtered,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := fs.FilteredStatus(c.inResult, c.inComp, c.inLog)
			require.NoError(t, err, "filtering shouldn't error")
			assert.Equal(t, c.want, got, "wrong status returned")
		})
	}
}
