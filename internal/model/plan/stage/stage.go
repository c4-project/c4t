// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package stage contains plan metadata describing which stages of a test cycle have occurred.
package stage

// Stage is the enumeration of stages.
//
// A stage generally corresponds to one of the 'act-tester-*' sub-programs, and represents a specific transformation
// of a plan file.
type Stage uint8

const (
	// Unknown is a sentinel value for unknown stages.
	Unknown Stage = iota

	// Plan is the required stage corresponding to selecting an input corpus and compiler set for future testing.
	Plan

	// Fuzz is the optional stage corresponding to mutating an input corpus.
	Fuzz

	// Lift is the required stage corresponding to generating test harnesses and build recipes for a corpus.
	Lift

	// Invoke is the required stage corresponding to running a plan against its machine node.
	Invoke

	// Analyse is the optional stage corresponding to post-processing an invoked plan.
	// Unlike other stages, it isn't logged in the plan file, and can be repeated.
	Analyse
)

//go:generate stringer -type Stage
