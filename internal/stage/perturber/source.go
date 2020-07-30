// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package perturber

import (
	"errors"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
)

var (
	ErrCListerNil    = errors.New("compiler lister nil")
	ErrCInspectorNil = errors.New("compiler inspector nil")
)

// Source contains all of the various sources for a Planner's information.
type Source struct {
	// CLister is the compiler lister.
	CLister CompilerLister

	// CInspector is the compiler [optimisation level] inspector.
	CInspector compiler.Inspector
}

// Check makes sure that all of this source's components are present and accounted for.
func (s *Source) Check() error {
	if s.CLister == nil {
		return ErrCListerNil
	}
	if s.CInspector == nil {
		return ErrCInspectorNil
	}
	return nil
}
