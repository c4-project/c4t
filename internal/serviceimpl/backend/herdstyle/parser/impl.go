// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package parser

import "github.com/c4-project/c4t/internal/subject/obs"

// Impl describes the parser functionality that differs between Herdtools-style backends.
type Impl interface {
	// ParseStateCount parses a potential state-count line whose raw fields are fields.
	// If the line is a state count, it returns (k, true, nil) where k is the state count;
	// if the line is to be skipped, it returns (_, false, nil); else, an error.
	ParseStateCount(fields []string) (k uint64, ok bool, err error)

	// ParseStateLine parses the state line whose raw fields are fields.
	ParseStateLine(tt TestType, fields []string) (*StateLine, error)

	// ParsePreTestLine extracts any flags implied by the pre-'Test' line whose raw fields are fields.
	//
	// Most implementations can safely pull this to (0, nil).
	// An example of an exception is Rmem, where the presence of the phrase 'PARTIAL RESULTS' before the 'Test' line
	// implies the 'partial' flag.
	ParsePreTestLine(fields []string) (obs.Flag, error)
}
