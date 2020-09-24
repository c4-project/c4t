// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject

import "errors"

var (
	// ErrMissingCompilation occurs on requests for compile results for a compiler that do not have them.
	ErrMissingCompilation = errors.New("no such compilation")

	// ErrDuplicateCompile occurs when one tries to insert a compile result that already exists.
	ErrDuplicateCompile = errors.New("duplicate compile result")

	// ErrDuplicateRecipe occurs when one tries to insert a recipe that already exists.
	ErrDuplicateRecipe = errors.New("duplicate recipe")

	// ErrDuplicateRun occurs when one tries to insert a run that already exists.
	ErrDuplicateRun = errors.New("duplicate run")

	// ErrMissingCompile occurs on requests for compile results for a compiler that do not have them.
	ErrMissingCompile = errors.New("no such compile result")

	// ErrMissingRecipe occurs on requests for recipe paths for an arch that do not have them.
	ErrMissingRecipe = errors.New("no such recipe")

	// ErrMissingRun occurs on requests for runs for a compiler that do not have them.
	ErrMissingRun = errors.New("no such run")

	// ErrNoBestLitmus occurs when asking for a BestLitmus() on a test with no valid Litmus file paths.
	ErrNoBestLitmus = errors.New("no valid litmus file for this subject")
)
