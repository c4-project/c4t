// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"errors"

	backend2 "github.com/MattWindsor91/c4t/internal/model/service/backend"
)

var (
	ErrBProbeNil  = errors.New("backend finder nil")
	ErrCListerNil = errors.New("compiler lister nil")
	ErrSProbeNil  = errors.New("subject prober nil")
)

// Source contains all of the various sources for a Planner's information.
type Source struct {
	// BProbe is the backend finder.
	BProbe backend2.Finder

	// CLister is the compiler lister.
	CLister CompilerLister

	// SProbe is the subject prober.
	SProbe SubjectProber
}

// Check makes sure that all of this source's components are present and accounted for.
func (s *Source) Check() error {
	if s.BProbe == nil {
		return ErrBProbeNil
	}
	if s.CLister == nil {
		return ErrCListerNil
	}
	if s.SProbe == nil {
		return ErrSProbeNil
	}
	return nil
}
