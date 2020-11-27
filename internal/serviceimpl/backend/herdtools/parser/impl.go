// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package parser

// Impl describes the parser functionality that differs between Herdtools-style backends.
type Impl interface {
	// ParseStateCount parses the state-count line whose raw fields are fields.
	ParseStateCount(fields []string) (uint64, error)

	// ParseStateLine parses the state line whose raw fields are fields.
	ParseStateLine(tt TestType, fields []string) (*StateLine, error)
}
