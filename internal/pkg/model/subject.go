package model

import (
	"errors"
	"fmt"
	"strings"
)

var ErrMissingHarness = errors.New("no such harness")

// Subject represents a single test subject in a corpus.
type Subject struct {
	// Name is the name of this subject.
	Name string `toml:"name"`

	// Litmus is the path to this subject's current Litmus file.
	Litmus string `toml:"name,omitempty"`

	// OrigLitmus is the path to this subject's original Litmus file.
	// If empty, then Litmus is the original file.
	OrigLitmus string `toml:"orig_litmus,omitempty"`

	// TracePath is the path to this subject's fuzzer trace file.
	// If empty, this subject hasn't been fuzzed by act-tester-fuzz.
	TracePath string `toml:"trace_path,omitempty"`

	// HarnessPaths contains the paths of every file in this subject's test harness.
	// Each maps from a string of the form 'machine:emits', where machine and emits are ACT IDs.
	// If nil, this subject hasn't had a harness generated.
	HarnessPaths map[string][]string `toml:"harness_paths,omitempty"`
}

// HarnessPath gets the harness path for the given machine and emits IDs.
func (s Subject) HarnessPath(machine, emits Id) ([]string, error) {
	key := strings.Join([]string{machine.String(), emits.String()}, ":")
	slice, ok := s.HarnessPaths[key]
	if !ok {
		return nil, fmt.Errorf("%w: machine=%q, emits=%q", ErrMissingHarness, machine, emits)
	}
	return slice, nil
}
