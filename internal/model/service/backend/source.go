// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package backend

// Source is the enumeration of possible input types to a lift job.
type Source uint8

const (
	// LiftUnknown states that the lifting source is unknown.
	LiftUnknown Source = iota
	// LiftLitmus states that the backend takes Litmus tests.
	LiftLitmus
)

//go:generate stringer -type Source
