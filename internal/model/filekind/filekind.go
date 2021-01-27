// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package filekind contains types for dealing with the various 'kinds' of file present in a subject.
package filekind

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
)

// Kind is the bitflag enumeration of file kinds.
type Kind uint8

const (
	// Other states that the kind of this file doesn't fit in any of the above categorisations.
	Other Kind = 1 << iota
	// Litmus states that this file is a litmus test.
	Litmus
	// Bin states that this file is a binary.
	Bin
	// Log states that this file is a compile log.
	Log
	// Trace states that this file is a fuzzer trace.
	Trace
	// CSrc states that this file is C source code (.c).
	CSrc
	// CHeader states that this file is a C header (.h).
	CHeader
	// Reserved for future use
	reserved

	// C is shorthand for CSrc|CHeader.
	C = CSrc | CHeader

	// Any is a suggestive alias for both Loc and Kind saturation.
	Any = math.MaxUint8

	strOther    = "other"
	strLitmus   = "litmus"
	strBin      = "bin"
	strLog      = "log"
	strTrace    = "trace"
	strCSrc     = "c/src"
	strCHeader  = "c/header"
	strC        = "c"
	strReserved = "RESERVED"
	sep         = "|"
)

// ErrBadKind occurs if we try to convert a kind from a string that doesn't match any known kind string.
var ErrBadKind = errors.New("unknown kind")

// Strings produces a stringified representation of each flag enabled in this filekind.
func (k Kind) Strings() []string {
	var sb []string
	add := func(j Kind, s string) bool {
		if !j.Matches(k) {
			return false
		}
		sb = append(sb, s)
		return true
	}

	add(Other, strOther)
	add(Litmus, strLitmus)
	add(Bin, strBin)
	add(Log, strLog)
	add(Trace, strTrace)
	add(reserved, strReserved)

	if !add(C, strC) {
		add(CHeader, strCHeader)
		add(CSrc, strCSrc)
	}

	return sb
}

// String produces a stringified representation of a filekind.
func (k Kind) String() string {
	if k == 0 {
		return "(none)"
	}
	return strings.Join(k.Strings(), sep)
}

// KindFromString tries to map str to a filekind.
func KindFromString(str string) (Kind, error) {
	switch strings.ToLower(str) {
	case strOther:
		return Other, nil
	case strLitmus:
		return Litmus, nil
	case strBin:
		return Bin, nil
	case strLog:
		return Log, nil
	case strTrace:
		return Trace, nil
	case strCHeader:
		return CHeader, nil
	case strCSrc:
		return CSrc, nil
	// Composite cases
	case strC:
		return C, nil
	default:
		return Other, fmt.Errorf("%w: %q", ErrBadKind, str)
	}
}

// KindFromStrings returns the union of the kinds corresponding to each string in strs.
func KindFromStrings(strs ...string) (Kind, error) {
	var k Kind
	for _, str := range strs {
		kd, err := KindFromString(str)
		if err != nil {
			return kd, err
		}
		k |= kd
	}
	return k, nil
}

// Matches checks whether this kind is included in pat.
func (k Kind) Matches(pat Kind) bool {
	return k&pat == k
}

// ArchivePerm gets the idealised Unix permission set for archiving a file of this kind.
func (k Kind) ArchivePerm() int64 {
	if k.Matches(Bin) {
		return 0755
	}
	return 0644
}

// MarshalJSON marshals a filekind to JSON using its string-list form.
func (k Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.Strings())
}

// UnmarshalJSON unmarshals an op from JSON using its string-list form.
func (k *Kind) UnmarshalJSON(bytes []byte) error {
	var (
		is  []string
		err error
	)
	if err = json.Unmarshal(bytes, &is); err != nil {
		return err
	}
	*k, err = KindFromStrings(is...)
	return err
}
