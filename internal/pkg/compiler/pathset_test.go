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

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// ExamplePathset_Dirs is a testable example for Dirs.
func ExamplePathset_Dirs() {
	ps := Pathset{DirBins: "bins", DirLogs: "logs"}
	for _, p := range ps.Dirs(model.IDFromString("foo"), model.IDFromString("bar.baz")) {
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

// TestPathset_OnCompiler tests that SubjectPaths produces sensible paths.
func TestPathset_OnCompiler(t *testing.T) {
	ps := Pathset{
		DirBins: "bins",
		DirLogs: "logs",
	}
	cid := model.IDFromString("foo.bar.baz")
	sps := ps.SubjectPaths(SubjectCompile{
		Name:       "yeet",
		CompilerID: cid,
	})

	wantb := path.Join("bins", "foo", "bar", "baz", "yeet")
	if sps.Bin != wantb {
		t.Errorf("bin for %s= %s, want %s", cid.String(), sps.Bin, wantb)
	}

	wantl := path.Join("logs", "foo", "bar", "baz", "yeet")
	if sps.Log != wantl {
		t.Errorf("logs for %s= %s, want %s", cid.String(), sps.Log, wantl)
	}
}
