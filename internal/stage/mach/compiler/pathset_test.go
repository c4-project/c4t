// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"fmt"
	"path"
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// ExamplePathset_Dirs is a testable example for Dirs.
func ExamplePathset_Dirs() {
	ps := Pathset{DirBins: "bins", DirLogs: "logs"}
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
	ps := NewPathset("foo")
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
	ps := Pathset{
		DirBins: "bins",
		DirLogs: "logs",
	}
	cid := id.FromString("foo.bar.baz")
	sps := ps.SubjectPaths(SubjectCompile{
		Name:       "yeet",
		CompilerID: cid,
	})

	wantb := path.Join("bins", "foo", "bar", "baz", "yeet")
	assert.Equal(t, wantb, sps.Bin, "bin on SubjectPaths not as expected")

	wantl := path.Join("logs", "foo", "bar", "baz", "yeet")
	assert.Equal(t, wantl, sps.Log, "log on SubjectPaths not as expected")
}
