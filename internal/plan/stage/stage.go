// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package stage contains plan metadata describing which stages of a test cycle have occurred.
package stage

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Stage is the enumeration of stages.
//
// A stage generally corresponds to one of the 'c4t-*' sub-programs, and represents a specific transformation
// of a plan file.
type Stage uint8

const (
	// Unknown is a sentinel value for unknown stages.
	Unknown Stage = iota

	// Plan is the required stage corresponding to selecting an input corpus and compiler set for future testing.
	Plan

	// Perturb is the optional stage corresponding to randomising and sampling a preceding plan.
	Perturb

	// Fuzz is the optional stage corresponding to mutating an input corpus.
	Fuzz

	// Lift is the required stage corresponding to generating test harnesses and build recipes for a corpus.
	Lift

	// Invoke is the required stage corresponding to compiling and running a plan.
	Invoke

	// Mach is a sub-stage of Invoke, corresponding to operations residing on the machine node.
	Mach

	// Compile is a sub-stage of Mach, corresponding to compiling the recipes in a plan.
	Compile

	// Run is a sub-stage of Mach, corresponding to running the compiled binaries in a plan.
	Run

	// Analyse is the optional stage corresponding to post-processing an invoked plan.
	// Unlike other stages, it isn't logged in the plan file, and can be repeated.
	Analyse

	// SetCompiler is the stage corresponding to manually setting a compiler.
	SetCompiler

	// Last points to the last stage in the enumeration.
	Last = SetCompiler
)

//go:generate stringer -type Stage

// FromString tries to convert a string into a Stage.
func FromString(s string) (Stage, error) {
	for i := Unknown; i <= Last; i++ {
		if strings.EqualFold(s, i.String()) {
			return i, nil
		}
	}
	return Unknown, fmt.Errorf("unknown Stage: %q", s)
}

// MarshalJSON marshals a stage to JSON using its string form.
func (i Stage) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON unmarshals a stage from JSON using its string form.
func (i *Stage) UnmarshalJSON(bytes []byte) error {
	var (
		is  string
		err error
	)
	if err = json.Unmarshal(bytes, &is); err != nil {
		return err
	}
	*i, err = FromString(is)
	return err
}
