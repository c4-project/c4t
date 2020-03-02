// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ObsFlag is the type of observation flags.
type ObsFlag int

const (
	// ObsSat represents a satisfying observation.
	ObsSat ObsFlag = 1 << iota
	// ObsUnat represents an unsatisfying observation.
	ObsUnsat
	// ObsUndef represents an undefined-behaviour observation.
	ObsUndef
)

var (
	// ErrBadObsFlag occurs when we read an unknown observation flag.
	ErrBadObsFlag = errors.New("bad observation flag")

	// ObsFlagNames maps the string representation of each observation flag to its flag value.
	ObsFlagNames = map[string]ObsFlag{
		"sat":   ObsSat,
		"unsat": ObsUnsat,
		"undef": ObsUndef,
	}
)

// Has checks to see if f is present in this flagset.
func (o ObsFlag) Has(f ObsFlag) bool {
	return o&f != ObsFlag(0)
}

// Strings expands this ObsFlag into string equivalents for each set flag.
func (o ObsFlag) Strings() []string {
	strs := make([]string, 0, 3)
	for str, f := range ObsFlagNames {
		if o.Has(f) {
			strs = append(strs, str)
		}
	}
	sort.Strings(strs)
	return strs
}

// ObsFlagOfStrings reconstitutes an observation flag given a representation as a list strs of strings.
func ObsFlagOfStrings(strs ...string) (ObsFlag, error) {
	var o ObsFlag
	for _, s := range strs {
		f, ok := ObsFlagNames[s]
		if !ok {
			return o, fmt.Errorf("%w: %s", ErrBadObsFlag, s)
		}
		o |= f
	}
	return o, nil
}

// MarshalText marshals an observation flag as a space-delimited string list.
func (o ObsFlag) MarshalText() ([]byte, error) {
	return []byte(strings.Join(o.Strings(), " ")), nil
}

// UnmarshalText unmarshals an observation flag list from bs by interpreting it as a string list.
func (o *ObsFlag) UnmarshalText(bs []byte) error {
	strs := strings.Fields(string(bs))

	var err error
	*o, err = ObsFlagOfStrings(strs...)
	return err
}

// An observed state.
type ObsState map[string]string

// Obs represents an observation in ACT's JSON-based format.
type Obs struct {
	// Flags contains any flags that are active on Obs.
	Flags ObsFlag `json:"flags,omitempty" toml:"flags,omitzero"`

	// CounterExamples lists all states that passed validation.
	CounterExamples []ObsState `json:"counter_examples,omitempty" toml:"counter_examples,omitempty"`

	// Witnesses lists all states that passed validation.
	Witnesses []ObsState `json:"witnesses,omitempty" toml:"witnesses,omitempty"`

	// States lists all observed states.
	States []ObsState `json:"states" toml:"states,omitempty"`
}

// Sat gets whether the observation satisfies its validation.
func (o *Obs) Sat() bool {
	return o.Flags.Has(ObsSat)
}

// Unsat gets whether the observation does not satisfy its validation.
func (o *Obs) Unsat() bool {
	return o.Flags.Has(ObsUnsat)
}
