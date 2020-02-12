package model

import (
	"errors"
	"fmt"
	"strings"
)

// ErrMissingHarness occurs on requests for harness paths for a machine/arch that do not have them.
var ErrMissingHarness = errors.New("no such harness")

// Subject represents a single test subject in a corpus.
type Subject struct {
	// Name is the name of this subject.
	Name string `toml:"name"`

	// Litmus is the path to this subject's current Litmus file.
	Litmus string `toml:"litmus,omitempty"`

	// OrigLitmus is the path to this subject's original Litmus file.
	// If empty, then Litmus is the original file.
	OrigLitmus string `toml:"orig_litmus,omitempty"`

	// TracePath is the path to this subject's fuzzer trace file.
	// If empty, this subject hasn't been fuzzed by act-tester-fuzz.
	TracePath string `toml:"trace_path,omitempty"`

	// Harnesses contains information about this subject's test harnesses.
	// It maps from a string of the form 'machine:arch', where machine and arch are ACT IDs.
	// If nil, this subject hasn't had a harness generated.
	Harnesses map[string]Harness `toml:"harnesses,omitempty"`
}

// Harness gets the harness for the given machine and arch IDs.
func (s *Subject) Harness(machine, arch ID) (Harness, error) {
	h, ok := s.Harnesses[harnessKey(machine, arch)]
	if !ok {
		return Harness{}, fmt.Errorf("%w: machine=%q, arch=%q", ErrMissingHarness, machine, arch)
	}
	return h, nil
}

// AddHarness sets the harness information for machine and arch to h in this subject.
func (s *Subject) AddHarness(machine, arch ID, h Harness) {
	if s.Harnesses == nil {
		s.Harnesses = make(map[string]Harness)
	}
	s.Harnesses[harnessKey(machine, arch)] = h
}

// harnessKey gets the harness-path key for a given machine and arch ID.
func harnessKey(machine, arch ID) string {
	return strings.Join([]string{machine.String(), arch.String()}, ":")
}
