// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend

// Source is the enumeration of possible input types to a lift job.
type Source uint8

const (
	// LiftUnknown states that the lifting source is unknown.
	LiftUnknown Source = iota
	// LiftLitmus states that the backend takes Litmus tests.
	// The backend may further specify particular architectures it can handle.
	LiftLitmus
)

//go:generate stringer -type Source
