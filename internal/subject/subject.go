// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package subject contains types and functions for dealing with test subject records.
// Such subjects generally live in a plan; the separate package exists to accommodate the large amount of subject
// specific types and functions in relation to the other parts of a test plan.

package subject

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// Normalise represents a single test subject in a corpus.
type Subject struct {
	// Fuzz is the fuzzer output for this subject, if it has been fuzzed.
	Fuzz *Fuzz `toml:"fuzz,omitempty" json:"fuzz,omitempty"`

	// Source refers to the original litmus test for this subject.
	Source litmus.Litmus `toml:"source,omitempty" json:"source,omitempty"`

	// Compiles contains information about this subject's compilation attempts.
	// It maps from the string form of each compiler's ID.
	// If nil, this subject hasn't had any compilations.
	Compiles map[string]compilation.CompileResult `toml:"compiles,omitempty" json:"compiles,omitempty"`

	// Recipes contains information about this subject's lifted test recipes.
	// It maps the string form of each recipe's target architecture's ID.
	// If nil, this subject hasn't had a recipe generated.
	Recipes map[string]recipe.Recipe `toml:"recipes,omitempty" json:"recipes,omitempty"`

	// Runs contains information about this subject's runs so far.
	// It maps from the string form of each compiler's ID.
	// If nil, this subject hasn't had any runs.
	Runs map[string]compilation.RunResult `toml:"runs,omitempty" json:"runs,omitempty"`
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

// Note that all of these maps work in basically the same way; their being separate and duplicated is just a
// consequence of Go not (yet) having generics.

// CompileResult gets the compilation result for the compiler ID cid.
func (s *Subject) CompileResult(cid id.ID) (compilation.CompileResult, error) {
	key := cid.String()
	c, ok := s.Compiles[key]
	if !ok {
		return compilation.CompileResult{}, fmt.Errorf("%w: compiler=%q", ErrMissingCompile, key)
	}
	return c, nil
}

// AddCompileResult sets the compilation information for compiler ID cid to c in this subject.
// It fails if there already _is_ a compilation.
func (s *Subject) AddCompileResult(cid id.ID, c compilation.CompileResult) error {
	s.ensureCompileMap()
	key := cid.String()
	if _, ok := s.Compiles[key]; ok {
		return fmt.Errorf("%w: compiler=%q", ErrDuplicateCompile, key)
	}
	s.Compiles[key] = c
	return nil
}

// ensureCompileMap makes sure this subject has a compile result map.
func (s *Subject) ensureCompileMap() {
	if s.Compiles == nil {
		s.Compiles = make(map[string]compilation.CompileResult)
	}
}

// Recipe gets the recipe for the architecture with id arch.
func (s *Subject) Recipe(arch id.ID) (recipe.Recipe, error) {
	key := arch.String()
	h, ok := s.Recipes[key]
	if !ok {
		return recipe.Recipe{}, fmt.Errorf("%w: arch=%q", ErrMissingRecipe, key)
	}
	return h, nil
}

// AddRecipe sets the recipe information for arch to r in this subject.
// It fails if there already _is_ a recipe for arch.
func (s *Subject) AddRecipe(arch id.ID, r recipe.Recipe) error {
	s.ensureRecipeMap()
	key := arch.String()
	if _, ok := s.Recipes[key]; ok {
		return fmt.Errorf("%w: arch=%q", ErrDuplicateRecipe, key)
	}
	s.Recipes[key] = r
	return nil
}

// ensureRecipeMap makes sure this subject has a recipe map.
func (s *Subject) ensureRecipeMap() {
	if s.Recipes == nil {
		s.Recipes = make(map[string]recipe.Recipe)
	}
}

// RunOf gets the run for the compiler with id cid.
func (s *Subject) RunOf(cid id.ID) (compilation.RunResult, error) {
	key := cid.String()
	h, ok := s.Runs[key]
	if !ok {
		return compilation.RunResult{}, fmt.Errorf("%w: compiler=%q", ErrMissingRun, key)
	}
	return h, nil
}

// AddRun sets the run information for cid to r in this subject.
// It fails if there already _is_ a run for cid.
func (s *Subject) AddRun(cid id.ID, r compilation.RunResult) error {
	s.ensureRunMap()
	key := cid.String()
	if _, ok := s.Runs[key]; ok {
		return fmt.Errorf("%w: compiler=%q", ErrDuplicateRun, key)
	}
	s.Runs[key] = r
	return nil
}

// ensureRunMap makes sure this subject has a run map.
func (s *Subject) ensureRunMap() {
	if s.Runs == nil {
		s.Runs = make(map[string]compilation.RunResult)
	}
}
