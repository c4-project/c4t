package runner

import (
	"path"
	"reflect"
	"sort"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// Test_Pathset_Dirs makes sure each of the expected paths appears in the pathset.
func TestPathset_Dirs(t *testing.T) {
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

// TestPathset_OnCompiler tests that OnCompiler produces sensible paths.
func TestPathset_OnCompiler(t *testing.T) {
	ps := Pathset{
		DirBins:     "bins",
		DirLogs:     "logs",
		DirFailures: "fails",
	}
	cid := model.IDFromString("foo.bar.baz")
	b, l := ps.OnCompiler(cid, "yeet")

	wantb := path.Join("bins", "foo", "bar", "baz", "yeet")
	if b != wantb {
		t.Errorf("bin for %s= %s, want %s", cid.String(), b, wantb)
	}

	wantl := path.Join("logs", "foo", "bar", "baz", "yeet")
	if l != wantl {
		t.Errorf("logs for %s= %s, want %s", cid.String(), l, wantl)
	}
}
