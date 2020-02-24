// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

// pathset contains the pre-computed paths used by the director.
type pathset struct {
	CurRun runPathset

	// DirLitmus is the directory to which litmus tests will be written.
	DirLitmus string

	// DirTrace is the directory to which traces will be written.
	DirTrace string
}

// runPathset contains the pre-computed paths used by a run of the director.
type runPathset struct {
	// DirFuzz is the fuzzing directory for this run.
	DirFuzz string

	// DirLift is the lifting directory for this run.
	DirLift string
}
