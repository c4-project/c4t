package model

import (
	"encoding/json"
	"errors"
	"fmt"
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
	return strs
}

// ObsFlagOfStrings reconstitutes an observation flag given a representation as a list strs of strings.
func ObsFlagOfStrings(strs []string) (ObsFlag, error) {
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

// MarshalJSON marshals an observation flag list as a string list.
func (o ObsFlag) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Strings())
}

// MarshalJSON unmarshals an observation flag list from bs by interpreting it as a string list.
func (o *ObsFlag) UnmarshalJSON(bs []byte) error {
	var strs []string
	if err := json.Unmarshal(bs, &strs); err != nil {
		return err
	}

	var err error
	*o, err = ObsFlagOfStrings(strs)
	return err
}

// An observed state.
type ObsState map[string]string

// Obs represents an observation in ACT's JSON-based format.
type Obs struct {
	// Flags contains any flags that are active on Obs.
	Flags ObsFlag `json:"flags"`

	// CounterExamples lists all states that passed validation.
	CounterExamples []ObsState `json:"counter_examples"`

	// Witnesses lists all states that passed validation.
	Witnesses []ObsState `json:"witnesses"`

	// States lists all observed states.
	States []ObsState `json:"states"`
}
