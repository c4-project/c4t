// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import "github.com/MattWindsor91/act-tester/internal/model/compiler"

// Source contains all of the various sources for a Planner's information.
type Source struct {
	// BProbe is the backend finder.
	BProbe BackendFinder

	// CLister is the compiler lister.
	CLister CompilerLister

	// CInspector is the compiler [optimisation level] inspector.
	CInspector compiler.Inspector

	// SProbe is the subject prober.
	SProbe SubjectProber
}
