// Package subject contains types and functions for dealing with test subject records.
// Such subjects generally live in a plan; the separate package exists to accommodate the large amount of subject
// specific types and functions in relation to the other parts of a test plan.

package subject

import (
	"errors"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

var (
	// ErrDuplicateCompile occurs when one tries to insert a compile result that already exists.
	ErrDuplicateCompile = errors.New("duplicate compile result")

	// ErrDuplicateHarness occurs when one tries to insert a harness that already exists.
	ErrDuplicateHarness = errors.New("duplicate harness")

	// ErrMissingCompile occurs on requests for compile results for a machine/compiler that do not have them.
	ErrMissingCompile = errors.New("no such compile result")

	// ErrMissingHarness occurs on requests for harness paths for a machine/arch that do not have them.
	ErrMissingHarness = errors.New("no such harness")

	// ErrNoBestLitmus occurs when asking for a BestLitmus() on a test with no valid Litmus file paths.
	ErrNoBestLitmus = errors.New("no valid litmus file for this subject")
)

// Subject represents a single test subject in a corpus.
type Subject struct {
	// Threads is the number of threads contained in this subject.
	Threads int `toml:"threads,omitzero"`

	// FuzzFileset is the fuzzing pathset for this subject, if it has been fuzzed.
	Fuzz *FuzzFileset `toml:"fuzz,omitempty"`

	// Litmus is the path to this subject's original Litmus file.
	Litmus string `toml:"litmus,omitempty"`

	// Compiles contains information about this subject's compilation attempts.
	// It maps from a string of the form 'machine:compiler', where machine and compiler are ACT IDs.
	// If nil, this subject hasn't had any compilations.
	Compiles map[string]CompileResult `toml:"compiles, omitempty"`

	// Harnesses contains information about this subject's test harnesses.
	// It maps from a string of the form 'machine:arch', where machine and arch are ACT IDs.
	// If nil, this subject hasn't had a harness generated.
	Harnesses map[string]Harness `toml:"harnesses,omitempty"`
}

// BestLitmus tries to get the 'best' litmus test path for further development.
//
// When there is a fuzzing record for this subject, the fuzz output is the best path.
// Otherwise, if there is a non-empty Litmus file for this subject, that file is the best path.
// Else, BestLitmus returns an error.
func (s *Subject) BestLitmus() (string, error) {
	switch {
	case s.Fuzz != nil && s.Fuzz.Litmus != "":
		return s.Fuzz.Litmus, nil
	case s.Litmus != "":
		return s.Litmus, nil
	default:
		return "", ErrNoBestLitmus
	}
}

// Note that the Compiles and Harnesses maps work in basically the same way; their being separate and duplicated is just a
// consequence of Go not (yet) having generics.

// CompileResult gets the compilation result for the machine-qualified compiler ID mcomp.
func (s *Subject) CompileResult(mcomp model.MachQualID) (CompileResult, error) {
	c, ok := s.Compiles[mcomp.String()]
	if !ok {
		return CompileResult{}, fmt.Errorf("%w: machine=%q, compiler=%q", ErrMissingCompile, mcomp.MachineID, mcomp.ID)
	}
	return c, nil
}

// AddCompileResult sets the compilation information for mcomp to c in this subject.
// It fails if there already _is_ a compilation.
func (s *Subject) AddCompileResult(mcomp model.MachQualID, c CompileResult) error {
	s.ensureCompileMap()
	key := mcomp.String()
	if _, ok := s.Compiles[key]; ok {
		return fmt.Errorf("%w: machine=%q, compiler=%q", ErrDuplicateCompile, mcomp.MachineID, mcomp.ID)
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

// Harness gets the harness for the machine-qualified arch ID.
func (s *Subject) Harness(march model.MachQualID) (Harness, error) {
	h, ok := s.Harnesses[march.String()]
	if !ok {
		return Harness{}, fmt.Errorf("%w: machine=%q, arch=%q", ErrMissingHarness, march.MachineID, march.ID)
	}
	return h, nil
}

// AddHarness sets the harness information for machine and arch to h in this subject.
// It fails if there already _is_ a harness.
func (s *Subject) AddHarness(march model.MachQualID, h Harness) error {
	s.ensureHarnessMap()
	key := march.String()
	if _, ok := s.Harnesses[key]; ok {
		return fmt.Errorf("%w: machine=%q, arch=%q", ErrDuplicateHarness, march.MachineID, march.ID)
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
