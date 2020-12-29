// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package normaliser

import (
	"path"

	"github.com/c4-project/c4t/internal/model/filekind"
)

// Map is the type of normalisation mappings.
type Map map[string]Entry

// Matching filters this map to only the files matching kind k and location l.
func (m Map) RenamesMatching(k filekind.Kind, l filekind.Loc) map[string]string {
	fs := make(map[string]string)
	for n, v := range m {
		if v.Kind.Matches(k) && v.Loc.Matches(l) {
			fs[n] = v.Original
		}
	}
	return fs
}

// Entry is a record in the normaliser's mappings.
// This exists mainly to make it possible to use a Normaliser to work out how to copy a plan to another host,
// but only copy selective subsets of files.
type Entry struct {
	// Original is the original path.
	Original string
	// Kind is the kind of path to which this mapping belongs.
	Kind filekind.Kind
	// Loc is an abstraction of the location of the path to which this mapping belongs.
	Loc filekind.Loc
}

// NewEntry constructs an entry with the filekind k, location l, and path constructed by segs.
func NewEntry(k filekind.Kind, l filekind.Loc, segs ...string) Entry {
	return Entry{
		Original: path.Join(segs...),
		Kind:     k,
		Loc:      l,
	}
}
