// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ObsFlag is the type of observation flags.
type Flag int

const (
	// Sat represents a satisfying observation.
	//
	// By default, this flag means that all states met the observation criteria.  If the Exist flag is also set,
	// it means that at least one state met the criteria.
	Sat Flag = 1 << iota
	// Unsat represents a satisfying observation.
	//
	// By default, this flag means that at least one state did not meet the observation criteria.  If the Exist flag is
	// also set, it means that no states met the criteria.
	Unsat
	// Undef represents an undefined-behaviour observation.
	Undef
	// Exist represents an existential observation (the default is for-all observations).
	Exist
)

var (
	// ErrBadFlag occurs when we read an unknown observation flag.
	ErrBadFlag = errors.New("bad observation flag")

	// FlagNames maps the string representation of each observation flag to its flag value.
	FlagNames = map[string]Flag{
		"sat":   Sat,
		"unsat": Unsat,
		"undef": Undef,
		"exist": Exist,
	}
)

// Has checks to see if f is present in this flagset.
func (o Flag) Has(f Flag) bool {
	return o&f != Flag(0)
}

// Strings expands this ObsFlag into string equivalents for each set flag.
func (o Flag) Strings() []string {
	strs := make([]string, 0, 3)
	for str, f := range FlagNames {
		if o.Has(f) {
			strs = append(strs, str)
		}
	}
	sort.Strings(strs)
	return strs
}

// FlagOfStrings reconstitutes an observation flag given a representation as a list strs of strings.
func FlagOfStrings(strs ...string) (Flag, error) {
	var o Flag
	for _, s := range strs {
		f, ok := FlagNames[s]
		if !ok {
			return o, fmt.Errorf("%w: %s", ErrBadFlag, s)
		}
		o |= f
	}
	return o, nil
}

// MarshalText marshals an observation flag as a space-delimited string list.
func (o Flag) MarshalText() ([]byte, error) {
	return []byte(strings.Join(o.Strings(), " ")), nil
}

// IsInteresting gets whether a flag represents an 'interesting' condition.
func (o Flag) IsInteresting() bool {
	return (o&(Undef) == Undef) || // Undefined flags are always interesting.
		(o&(Sat|Exist) == (Sat | Exist)) || // Satisfied flags are interesting if they are existential.
		(o&(Unsat|Exist) == Unsat) || // Unsatisfied flags are interesting if they are not existential.
		(o&(Sat|Unsat) == 0) // Flags that are neither sat nor unsat are interesting; they suggest something weird happened.
}

// IsSat gets whether a flag represents a satisfying observation.
func (o Flag) IsSat() bool {
	return o.Has(Sat)
}

// IsUnsat gets whether the observation does not satisfy its validation.
func (o Flag) IsUnsat() bool {
	return o.Has(Unsat)
}

// IsExistential gets whether the observation's validation was existential rather than universal.
func (o Flag) IsExistential() bool {
	return o.Has(Exist)
}

// UnmarshalText unmarshals an observation flag list from bs by interpreting it as a string list.
func (o *Flag) UnmarshalText(bs []byte) error {
	strs := strings.Fields(string(bs))

	var err error
	*o, err = FlagOfStrings(strs...)
	return err
}

// MarshalJSON marshals an observation flag list as a string list.
func (o Flag) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Strings())
}

// MarshalJSON unmarshals an observation flag list from bs by interpreting it as a string list.
func (o *Flag) UnmarshalJSON(bs []byte) error {
	var strs []string
	if err := json.Unmarshal(bs, &strs); err != nil {
		return err
	}

	var err error
	*o, err = FlagOfStrings(strs...)
	return err
}
