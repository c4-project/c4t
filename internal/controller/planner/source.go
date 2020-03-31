// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

// Source contains all of the various sources for a Planner's information.
type Source struct {
	// BProbe is the backend prober.
	BProbe BackendFinder

	// CProbe is the compiler prober.
	CProbe CompilerLister

	// SProbe is the subject prober.
	SProbe SubjectProber
}
