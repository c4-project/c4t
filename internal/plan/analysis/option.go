// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import "github.com/1set/gut/ystring"

// Option is the type of options to Analyse.
type Option func(*analyser) error

// Options applies each option in opts onto the analyser.
func Options(opts ...Option) Option {
	return func(a *analyser) error {
		for _, o := range opts {
			if err := o(a); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithWorkerCount sets the worker count of the analyser to nworkers.
func WithWorkerCount(nworkers int) Option {
	return func(a *analyser) error {
		a.nworkers = nworkers
		return nil
	}
}

// WithFilters appends the filters in fs to the filter set.
func WithFilters(fs FilterSet) Option {
	return func(a *analyser) error {
		a.filters = append(a.filters, fs...)
		return nil
	}
}

// WithFiltersFromFile appends the filters in the filter file at path to the filter set.
// If the path is empty, we don't add any filters.
func WithFiltersFromFile(path string) Option {
	return func(a *analyser) error {
		if ystring.IsBlank(path) {
			return nil
		}
		fs, err := LoadFilterSet(path)
		if err != nil {
			return err
		}
		return WithFilters(fs)(a)
	}
}
