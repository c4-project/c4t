// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package analysis_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/c4-project/c4t/internal/helper/testhelp"

	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/plan/analysis"
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

// TestLoadFilterSet tests LoadFilterSet on an example filter file that doesn't exist.
func TestLoadFilterSet_notFound(t *testing.T) {
	t.Parallel()

	_, err := analysis.LoadFilterSet(filepath.Join("testdata", "nonsuch.yaml"))
	testhelp.ExpectErrorIs(t, err, os.ErrNotExist, "loading nonexistent filter set")
}

// TestFilterSet_FilteredStatus tests FilterSet.FilteredStatus with several cases.
func TestFilterSet_FilteredStatus(t *testing.T) {
	t.Parallel()

	fs, err := analysis.LoadFilterSet(filepath.Join("testdata", "filters.yaml"))
	require.NoError(t, err, "loading filter set should not error")

	cases := map[string]struct {
		inComp     compiler.Instance
		inLog      string
		inStatus   status.Status
		fsOverride analysis.FilterSet
		want       status.Status
		err        error
	}{
		"no filtering": {
			inComp:   compiler.MockX86Gcc(),
			inLog:    "blep",
			inStatus: status.Ok,
			want:     status.Ok,
		},
		"filtering in middle of string": {
			inComp:   compiler.MockX86Gcc(),
			inLog:    "foo error: invalid memory model for ‘__atomic_exchange’ bar",
			inStatus: status.Ok,
			want:     status.Filtered,
		},
		"no filters": {
			fsOverride: analysis.FilterSet{},
			inComp:     compiler.MockX86Gcc(),
			inLog:      "foo error: invalid memory model for ‘__atomic_exchange’ bar",
			inStatus:   status.Ok,
			want:       status.Ok,
		},
		"filtering with a broken filter set": {
			fsOverride: analysis.FilterSet{
				{
					Style:        id.FromString("bad.*.glob.*"),
					ErrorPattern: "blep",
				},
			},
			inComp:   compiler.MockX86Gcc(),
			inLog:    "blep",
			inStatus: status.Ok,
			err:      id.ErrBadGlob,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := fs
			if c.fsOverride != nil {
				fs = c.fsOverride
			}

			got, err := fs.FilteredStatus(c.inStatus, c.inComp, c.inLog)
			if c.err != nil {
				testhelp.ExpectErrorIs(t, err, c.err, "FilteredStatus")
				return
			}
			require.NoError(t, err, "filtering shouldn't error")
			assert.Equal(t, c.want, got, "wrong status returned")
		})
	}
}
