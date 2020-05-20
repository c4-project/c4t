// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package filekind contains types for dealing with the various 'kinds' of file present in a subject.
package filekind

import "math"

// Kind is the bitflag enumeration of file kinds.
type Kind uint8

// Any is a suggestive alias for both Loc and Kind saturation.
const Any = math.MaxUint8

const (
	// Litmus states that this file is a litmus test.
	Litmus Kind = 1 << iota
	// Bin states that this file is a binary.
	Bin
	// CSrc states that this file is C source code (.c).
	CSrc
	// CHeader states that this file is a C header (.h).
	CHeader
	// Log states that this file is a compile log.
	Log
	// Trace states that this file is a fuzzer trace.
	Trace
	// Other states that the kind of this file doesn't fit in any of the above categorisations.
	Other

	// C is shorthand for CSrc|CHeader.
	C = CSrc | CHeader
)

// Matches checks whether this kind is included in pat.
func (k Kind) Matches(pat Kind) bool {
	return k&pat == k
}
