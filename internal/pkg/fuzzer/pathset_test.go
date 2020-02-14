package fuzzer

import (
	"reflect"
	"sort"
	"testing"
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
