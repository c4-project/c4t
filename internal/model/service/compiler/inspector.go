// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"errors"
	"fmt"

	"github.com/c4-project/c4t/internal/helper/stringhelp"

	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"
)

var (
	// ErrConfigNil occurs if we try to select optimisation levels for a nil compiler.
	ErrConfigNil = errors.New("can't select levels for nil compiler")
	// ErrNoSuchLevel occurs if a Selection enables an optimisation level that isn't available.
	ErrNoSuchLevel = errors.New("no such optimisation level")
)

// Inspector is the interface of types that support optimisation level lookup.
type Inspector interface {
	// DefaultOptLevels retrieves a set of optimisation levels that are enabled by default for compiler c.
	DefaultOptLevels(c *Compiler) (stringhelp.Set, error)
	// OptLevels retrieves a set of potential optimisation levels for compiler c.
	// This map shouldn't be modified, as it may be global.
	OptLevels(c *Compiler) (map[string]optlevel.Level, error)
	// DefaultMOpts retrieves a set of machine optimisation directives that are enabled by default for compiler c.
	DefaultMOpts(c *Compiler) (stringhelp.Set, error)

	// We don't request a list of all possible MOpts from compilers, as the list expands so rapidly that any such
	// list would be hideously out of date.
}

//go:generate mockery --name=Inspector

// SelectLevels selects from in the optimisation levels permitted by the configuration c.
func SelectLevels(in Inspector, c *Compiler) (map[string]optlevel.Level, error) {
	if c == nil {
		return nil, ErrConfigNil
	}

	all, err := in.OptLevels(c)
	if err != nil {
		return nil, err
	}
	dls, err := in.DefaultOptLevels(c)
	if err != nil {
		return nil, err
	}
	return filterLevels(chosenLevels(dls, c.Opt), all)
}

// SelectMOpts selects from in the machine optimisation profiles (-march, etc.) permitted by the configuration c.
func SelectMOpts(in Inspector, c *Compiler) (stringhelp.Set, error) {
	if c == nil {
		return nil, ErrConfigNil
	}

	dls, err := in.DefaultMOpts(c)
	if err != nil {
		return nil, err
	}
	return chosenLevels(dls, c.MOpt), nil
}

func chosenLevels(defaults stringhelp.Set, s *optlevel.Selection) stringhelp.Set {
	if s == nil {
		return defaults
	}
	return s.Override(defaults)
}

func filterLevels(choices stringhelp.Set, all map[string]optlevel.Level) (map[string]optlevel.Level, error) {
	chosen := make(map[string]optlevel.Level, len(choices))
	for c := range choices {
		var ok bool
		if chosen[c], ok = all[c]; !ok {
			return nil, fmt.Errorf("%w: %q", ErrNoSuchLevel, c)
		}
	}

	return chosen, nil
}
