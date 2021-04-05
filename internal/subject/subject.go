// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package subject contains types and functions for dealing with test subject records.
// Such subjects generally live in a plan; the separate package exists to accommodate the large amount of subject
// specific types and functions in relation to the other parts of a test plan.

package subject

import (
	"fmt"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/model/litmus"

	"github.com/c4-project/c4t/internal/model/recipe"

	"github.com/c4-project/c4t/internal/id"
)

// Normalise represents a single test subject in a corpus.
type Subject struct {
	// Fuzz is the fuzzer output for this subject, if it has been fuzzed.
	Fuzz *Fuzz `toml:"fuzz,omitempty" json:"fuzz,omitempty"`

	// Source refers to the original litmus test for this subject.
	Source litmus.Litmus `toml:"source,omitempty" json:"source,omitempty"`

	// Compilations contains information about this subject's compilations, organised by compiler ID.
	// If nil, this subject hasn't had any compilations.
	Compilations compilation.Map `toml:"compilations,omitempty" json:"compilations,omitempty"`

	// Recipes contains information about this subject's lifted test recipes.
	// If nil, this subject hasn't had any recipes generated.
	Recipes recipe.Map `toml:"recipes,omitempty" json:"recipes,omitempty"`
}

// BestLitmus tries to get the 'best' litmus test for further development.
//
// When there is a fuzzing record for this subject, the fuzz output is the best path.
// Otherwise, if there is a non-empty Litmus file for this subject, that file is the best path.
// Else, BestLitmus returns an error.
func (s *Subject) BestLitmus() (*litmus.Litmus, error) {
	switch {
	case s.HasFuzzFile():
		return &s.Fuzz.Litmus, nil
	case s.Source.HasPath():
		return &s.Source, nil
	default:
		return nil, ErrNoBestLitmus
	}
}

// HasFuzzFile gets whether this subject has a fuzzed testcase file.
func (s *Subject) HasFuzzFile() bool {
	return s.Fuzz != nil && s.Fuzz.Litmus.HasPath()
}

// Compilation gets the compilation information for the compiler ID cid.
func (s *Subject) Compilation(cid id.ID) (compilation.Compilation, error) {
	// TODO(@MattWindsor91): obsolete?
	c, ok := s.Compilations[cid]
	if !ok {
		return compilation.Compilation{}, fmt.Errorf("%w: compiler=%q", ErrMissingCompilation, cid)
	}
	return c, nil
}

// CompileResult gets the compilation result for the compiler ID cid.
func (s *Subject) CompileResult(cid id.ID) (*compilation.CompileResult, error) {
	c, err := s.Compilation(cid)
	if err != nil {
		return nil, err
	}
	if c.Compile == nil {
		return nil, fmt.Errorf("%w: compiler=%q", ErrMissingCompile, cid)
	}
	return c.Compile, err
}

// AddCompileResult sets the compilation information for compiler ID cid to c in this subject.
// It fails if there already _is_ a compilation.
func (s *Subject) AddCompileResult(cid id.ID, c compilation.CompileResult) error {
	return s.mapCompilation(cid, func(cc *compilation.Compilation) error {
		if cc.Compile != nil {
			return fmt.Errorf("%w: compiler=%q", ErrDuplicateCompile, cid)
		}
		cc.Compile = &c
		return nil
	})
}

// RunResult gets the run result for the compiler with id cid.
func (s *Subject) RunResult(cid id.ID) (*compilation.RunResult, error) {
	cc, err := s.Compilation(cid)
	if err != nil {
		return nil, err
	}
	if cc.Run == nil {
		return nil, fmt.Errorf("%w: compiler=%q", ErrMissingRun, cid)
	}
	return cc.Run, err
}

// AddRun sets the run information for cid to r in this subject.
// It fails if there already _is_ a run for cid.
func (s *Subject) AddRun(cid id.ID, r compilation.RunResult) error {
	return s.mapCompilation(cid, func(cc *compilation.Compilation) error {
		if cc.Run != nil {
			return fmt.Errorf("%w: compiler=%q", ErrDuplicateRun, cid)
		}
		cc.Run = &r
		return nil
	})
}

func (s *Subject) mapCompilation(cid id.ID, f func(cc *compilation.Compilation) error) error {
	s.ensureCompilationMap()
	// Deliberately taking the zero value if the compilation hasn't been seen yet.
	cc := s.Compilations[cid]
	if err := f(&cc); err != nil {
		return err
	}
	s.Compilations[cid] = cc
	return nil
}

// ensureCompilationMap makes sure this subject has a compile result map.
func (s *Subject) ensureCompilationMap() {
	if s.Compilations == nil {
		s.Compilations = make(compilation.Map)
	}
}

// Recipe gets the recipe for the architecture with id arch.
// It returns the ID of the recipe as well as the recipe contents.
func (s *Subject) Recipe(arch id.ID) (id.ID, recipe.Recipe, error) {
	// TODO(@MattWindsor91): do scoping here
	r, ok := s.Recipes[arch]
	if !ok {
		return id.ID{}, recipe.Recipe{}, fmt.Errorf("%w: arch=%q", ErrMissingRecipe, arch)
	}
	return arch, r, nil
}

// AddRecipe sets the recipe information for arch to r in this subject.
// It fails if there already _is_ a recipe for arch.
func (s *Subject) AddRecipe(arch id.ID, r recipe.Recipe) error {
	s.ensureRecipeMap()
	if _, ok := s.Recipes[arch]; ok {
		return fmt.Errorf("%w: arch=%q", ErrDuplicateRecipe, arch)
	}
	s.Recipes[arch] = r
	return nil
}

// ensureRecipeMap makes sure this subject has a recipe map.
func (s *Subject) ensureRecipeMap() {
	if s.Recipes == nil {
		s.Recipes = make(recipe.Map)
	}
}
