// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"errors"
	"fmt"
	"strings"
)

// ErrUnknownProfileKind is an error that occurs if we try to unmarshal an unknown profile kind.
var ErrUnknownProfileKind = errors.New("unknown profile kind")

// ProfileKind is the enumeration of kinds of coverage profile.
type ProfileKind uint8

const (
	// Known is a profile kind that tells the coverage generator to run a mutating fuzzer known to it.
	// At time of writing, there is only one such fuzzer (act-fuzz).
	Known ProfileKind = iota
	// Standalone is a profile kind that tells the coverage generator to run an external, stand-alone fuzzer.
	Standalone
	// LastProfileKind represents the last profile kind.
	LastProfileKind = Standalone
)

//go:generate stringer -type ProfileKind

// MarshalText marshals a profile kind to text by using its string representation.
func (i ProfileKind) MarshalText() (text []byte, err error) {
	return []byte(i.String()), nil
}

// UnmarshalText tries to unmarshal a profile kind from text.
func (i *ProfileKind) UnmarshalText(text []byte) error {
	s := string(text)
	for *i = Known; *i <= LastProfileKind; *i++ {
		if strings.EqualFold(i.String(), s) {
			return nil
		}
	}
	return fmt.Errorf("%w: %q", ErrUnknownProfileKind, s)
}
