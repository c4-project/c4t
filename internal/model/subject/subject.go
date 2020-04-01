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

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// Subject represents a single test subject in a corpus.
type Subject struct {
	// Threads is the number of threads contained in this subject.
	Threads int `toml:"threads,omitzero"`

	// Fuzz is the fuzzing pathset for this subject, if it has been fuzzed.
	Fuzz *Fuzz `toml:"fuzz,omitempty"`

	// Litmus is the (slashed) path to this subject's original Litmus file.
	Litmus string `toml:"litmus,omitempty"`

	// Compiles contains information about this subject's compilation attempts.
	// It maps from the string form of each compiler's ID.
	// If nil, this subject hasn't had any compilations.
	Compiles map[string]CompileResult `toml:"compiles, omitempty"`

	// Harnesses contains information about this subject's test harnesses.
	// It maps the string form of each harness's target architecture's ID.
	// If nil, this subject hasn't had a harness generated.
	Harnesses map[string]Harness `toml:"harnesses,omitempty"`

	// Runs contains information about this subject's runs so far.
	// It maps from the string form of each compiler's ID.
	// If nil, this subject hasn't had any runs.
	Runs map[string]RunResult `toml:"runs, omitempty"`
}

// BestLitmus tries to get the 'best' litmus test path for further development.
//
// When there is a fuzzing record for this subject, the fuzz output is the best path.
// Otherwise, if there is a non-empty Litmus file for this subject, that file is the best path.
// Else, BestLitmus returns an error.
func (s *Subject) BestLitmus() (string, error) {
	switch {
	case s.Fuzz != nil && s.Fuzz.Files.Litmus != "":
		return s.Fuzz.Files.Litmus, nil
	case s.Litmus != "":
		return s.Litmus, nil
	default:
		return "", ErrNoBestLitmus
	}
}

// Note that all of these maps work in basically the same way; their being separate and duplicated is just a
// consequence of Go not (yet) having generics.

// CompileResult gets the compilation result for the compiler ID cid.
func (s *Subject) CompileResult(cid id.ID) (CompileResult, error) {
	key := cid.String()
	c, ok := s.Compiles[key]
	if !ok {
		return CompileResult{}, fmt.Errorf("%w: compiler=%q", ErrMissingCompile, key)
	}
	return c, nil
}

// AddCompileResult sets the compilation information for compiler ID cid to c in this subject.
// It fails if there already _is_ a compilation.
func (s *Subject) AddCompileResult(cid id.ID, c CompileResult) error {
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
		s.Compiles = make(map[string]CompileResult)
	}
}

// Harness gets the harness for the architecture with id arch.
func (s *Subject) Harness(arch id.ID) (Harness, error) {
	key := arch.String()
	h, ok := s.Harnesses[key]
	if !ok {
		return Harness{}, fmt.Errorf("%w: arch=%q", ErrMissingHarness, key)
	}
	return h, nil
}

// AddHarness sets the harness information for arch to h in this subject.
// It fails if there already _is_ a harness for arch.
func (s *Subject) AddHarness(arch id.ID, h Harness) error {
	s.ensureHarnessMap()
	key := arch.String()
	if _, ok := s.Harnesses[key]; ok {
		return fmt.Errorf("%w: arch=%q", ErrDuplicateHarness, key)
	}
	s.Harnesses[key] = h
	return nil
}

// ensureHarnessMap makes sure this subject has a harness map.
func (s *Subject) ensureHarnessMap() {
	if s.Harnesses == nil {
		s.Harnesses = make(map[string]Harness)
	}
}

// RunOf gets the run for the compiler with id cid.
func (s *Subject) RunOf(cid id.ID) (RunResult, error) {
	key := cid.String()
	h, ok := s.Runs[key]
	if !ok {
		return RunResult{}, fmt.Errorf("%w: compiler=%q", ErrMissingRun, key)
	}
	return h, nil
}

// AddRun sets the run information for cid to r in this subject.
// It fails if there already _is_ a run for cid.
func (s *Subject) AddRun(cid id.ID, r RunResult) error {
	s.ensureRunMap()
	key := cid.String()
	if _, ok := s.Runs[key]; ok {
		return fmt.Errorf("%w: compiler=%q", ErrDuplicateRun, key)
	}
	s.Runs[key] = r
	return nil
}

// ensureHarnessMap makes sure this subject has a harness map.
func (s *Subject) ensureRunMap() {
	if s.Runs == nil {
		s.Runs = make(map[string]RunResult)
	}
}
