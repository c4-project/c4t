// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject

import "errors"

var (
	// ErrDuplicateCompile occurs when one tries to insert a compile result that already exists.
	ErrDuplicateCompile = errors.New("duplicate compile result")

	// ErrDuplicateHarness occurs when one tries to insert a harness that already exists.
	ErrDuplicateHarness = errors.New("duplicate harness")

	// ErrDuplicateRun occurs when one tries to insert a run that already exists.
	ErrDuplicateRun = errors.New("duplicate run")

	// ErrMissingCompile occurs on requests for compile results for a compiler that do not have them.
	ErrMissingCompile = errors.New("no such compile result")

	// ErrMissingHarness occurs on requests for harness paths for an arch that do not have them.
	ErrMissingHarness = errors.New("no such harness")

	// ErrMissingRun occurs on requests for runs for a compiler that do not have them.
	ErrMissingRun = errors.New("no such run")

	// ErrNoBestLitmus occurs when asking for a BestLitmus() on a test with no valid Litmus file paths.
	ErrNoBestLitmus = errors.New("no valid litmus file for this subject")
)
