// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"errors"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/compiler/optlevel"
)

var (
	// ErrConfigNil occurs if we try to select optimisation levels for a nil compiler.
	ErrConfigNil = errors.New("can't select levels for nil compiler")
	// ErrNoSuchLevel occurs if a Selection enables an optimisation level that isn't available.
	ErrNoSuchLevel = errors.New("no such optimisation level")
)

// Inspector is the interface of types that support optimisation level lookup.
type Inspector interface {
	// DefaultLevels retrieves a set of optimisation levels that are enabled by default for compiler c.
	DefaultLevels(c *Config) (map[string]struct{}, error)
	// Levels retrieves a set of potential optimisation levels for compiler c.
	// This map shouldn't be modified, as it may be global.
	Levels(c *Config) (map[string]optlevel.Level, error)
}

// SelectLevels selects from in the optimisation levels permitted by the configuration c.
func SelectLevels(in Inspector, c *Config) (map[string]optlevel.Level, error) {
	if c == nil {
		return nil, ErrConfigNil
	}

	all, err := in.Levels(c)
	if err != nil {
		return nil, err
	}
	dls, err := in.DefaultLevels(c)
	if err != nil {
		return nil, err
	}
	return filterLevels(chosenLevels(dls, c.Opt), all)
}

func chosenLevels(defaults map[string]struct{}, s *optlevel.Selection) map[string]struct{} {
	if s == nil {
		return defaults
	}
	return s.Override(defaults)
}

func filterLevels(choices map[string]struct{}, all map[string]optlevel.Level) (map[string]optlevel.Level, error) {
	chosen := make(map[string]optlevel.Level, len(choices))
	for c := range choices {
		var ok bool
		if chosen[c], ok = all[c]; !ok {
			return nil, fmt.Errorf("%w: %q", ErrNoSuchLevel, c)
		}
	}

	return chosen, nil
}
