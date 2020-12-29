// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"errors"
	"io"
	"os"
	"regexp"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/c4-project/c4t/internal/helper/errhelp"
	"github.com/c4-project/c4t/internal/model/id"
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
	// compiledPattern is the compiled version of ErrorPattern.
	compiledPattern *regexp.Regexp
}

// FilterSet is the type of sets of filter.
type FilterSet []Filter

// ReadFilterSet loads and compiles a filter set from the reader r.
func ReadFilterSet(r io.Reader) (FilterSet, error) {
	yd := yaml.NewDecoder(r)
	var fs FilterSet
	err := yd.Decode(&fs)
	if err != nil {
		return nil, err
	}
	return Compile(fs)
}

// Compile compiles the filter set fs.
func Compile(fs FilterSet) (FilterSet, error) {
	var err error
	for i := range fs {
		if fs[i].compiledPattern, err = regexp.Compile(fs[i].ErrorPattern); err != nil {
			return nil, err
		}
	}
	return fs, nil
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

// Filter returns true if, and only if, at least one filter in this set matches ci, and log.
func (f FilterSet) Filter(ci compiler.Configuration, log string) (bool, error) {
	for _, fl := range f {
		matched, err := fl.Filter(ci, log)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

// FilteredStatus returns s if FilterSet.Filter returns false over ci and log, or Filtered otherwise.
func (f FilterSet) FilteredStatus(s status.Status, ci compiler.Configuration, log string) (status.Status, error) {
	filtered, err := f.Filter(ci, log)
	if filtered {
		s = status.Filtered
	}
	return s, err
}

// Filter returns true if, and only if, this filter matches cm, ci, and log.
func (f Filter) Filter(ci compiler.Configuration, log string) (bool, error) {
	styleMatch, err := ci.Style.Matches(f.Style)
	if err != nil || !styleMatch {
		return false, err
	}
	// TODO(@MattWindsor91): compiler versions
	return f.filterCompilerLog(log)
}

func (f Filter) filterCompilerLog(log string) (bool, error) {
	if f.compiledPattern == nil {
		return false, errors.New("filter was not compiled")
	}
	return f.compiledPattern.MatchString(log), nil
}
