// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/c4-project/c4t/internal/stage/mach/compiler"

	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/id"
)

// ExamplePathset_Dirs is a testable example for Dirs.
func ExamplePathset_Dirs() {
	ps := compiler.Pathset{DirBins: "bins", DirLogs: "logs"}
	for _, p := range ps.Dirs(id.FromString("foo"), id.FromString("bar.baz")) {
		fmt.Println(p)
	}
	// Unordered output:
	// bins
	// bins/foo
	// bins/bar/baz
	// logs
	// logs/foo
	// logs/bar/baz
}

// Test_Pathset_Dirs_NoCompilers makes sure each of the expected paths appears in the pathset
// when no compilers are involved.
func TestPathset_Dirs_NoCompilers(t *testing.T) {
	ps := compiler.NewPathset("foo")
	dirs := ps.Dirs()
	sort.Strings(dirs)

	vps := reflect.ValueOf(*ps)
	tps := vps.Type()

	nf := vps.NumField()
	for i := 0; i < nf; i++ {
		dir := vps.Field(i).String()
		name := tps.Field(i).Name

		if r := sort.SearchStrings(dirs, dir); r < 0 || len(dirs) <= r {
			t.Errorf("missing %s (val %s) in dirs %v", name, dir, dirs)
		}
	}
}

// TestPathset_SubjectPaths tests that SubjectPaths produces sensible paths.
func TestPathset_SubjectPaths(t *testing.T) {
	ps := compiler.Pathset{
		DirBins: "bins",
		DirLogs: "logs",
	}
	cid := id.FromString("foo.bar.baz")
	sps := ps.SubjectPaths(compilation.Name{
		SubjectName: "yeet",
		CompilerID:  cid,
	})

	// Bin and Log are slashpaths, perhaps surprisingly.

	wantb := path.Join("bins", "foo", "bar", "baz", "yeet")
	assert.Equal(t, wantb, sps.Bin, "bin on SubjectPaths not as expected")

	wantl := path.Join("logs", "foo", "bar", "baz", "yeet")
	assert.Equal(t, wantl, sps.Log, "log on SubjectPaths not as expected")
}

func TestPathset_Prepare(t *testing.T) {
	td := t.TempDir()
	ps := compiler.NewPathset(td)

	compilers := []id.ID{id.FromString("gcc.4"), id.FromString("gcc.9"), id.FromString("icc")}

	err := ps.Prepare(compilers...)
	require.NoError(t, err, "preparing compile pathset in temp dir")

	for _, c := range compilers {
		cfs := ps.SubjectPaths(compilation.Name{
			SubjectName: "foo",
			CompilerID:  c,
		})

		// These will probably be the same directory, but there's no invariant to enforce that.
		assert.DirExists(t, filepath.Dir(filepath.Clean(cfs.Log)), "log directory should exist")
		assert.DirExists(t, filepath.Dir(filepath.Clean(cfs.Bin)), "bin directory should exist")
	}
}
