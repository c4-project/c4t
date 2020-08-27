// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"io"
	"os"

	"github.com/MattWindsor91/act-tester/internal/helper/errhelp"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"gopkg.in/yaml.v3"
)

// Filter is the type of filters.
//
// A filter, when triggered during analysis of a compilation, sets the 'filtered' flag on that compilation.
// This suppresses any non-OK status, replacing it with the 'filtered' status.
type Filter struct {
	// Style is a glob identifier that selects a particular compiler style.
	Style id.ID `yaml:"style"`
	// MajorVersionBelow is a lower bound on the major version of the compiler, if set to a positive number.
	MajorVersionBelow int `yaml:"major_version_below,omitempty"`
	// ErrorPattern is an uncompiled regexp that selects a particular phrase in a compiler error.
	ErrorPattern string `yaml:"error_pattern,omitempty"`
}

// FilterSet is the type of sets of filter.
type FilterSet []Filter

// ReadFilterSet loads a filter set from the reader r.
func ReadFilterSet(r io.Reader) (FilterSet, error) {
	yd := yaml.NewDecoder(r)
	var fs FilterSet
	err := yd.Decode(&fs)
	return fs, err
}

// LoadFilterSet loads a filter set from the filepath fpath.
func LoadFilterSet(fpath string) (FilterSet, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	fs, rerr := ReadFilterSet(f)
	cerr := f.Close()
	return fs, errhelp.FirstError(rerr, cerr)
}
