// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package rmem

import (
	"github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle/herd"
	"github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle/litmus"
	"github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle/parser"
	"github.com/c4-project/c4t/internal/subject/obs"
)

// ParseStateCount parses the state count in fields according to Rmem's syntax.
func (Rmem) ParseStateCount(fields []string) (uint64, bool, error) {
	// Rmem's syntax here is exactly the same as Herd's:
	return herd.Herd{}.ParseStateCount(fields)
}

func (Rmem) ParseStateLine(_ parser.TestType, fields []string) (*parser.StateLine, error) {
	// Rmem's state line syntax is similar to that of Litmus's, but with a few gotchas:
	//
	// - Asterisks always represent witnesses, not 'interesting' cases; this means we always parse as if the test type
	//   is 'allowed';
	// - State lines contain a 'via "XYZ"' line, which we need to scrub.
	return litmus.Litmus{}.ParseStateLine(parser.Allowed, stripVia(fields))
}

// ParsePreTestLine checks the pre-Test line in fields to check whether it states this is a partial observation.
func (Rmem) ParsePreTestLine(fields []string) (obs.Flag, error) {
	nf := len(fields)
	var f obs.Flag
	// This is preceded and succeeded by a large number of *s, but we don't check for those.
	if nf == 4 && fields[0] == "***" && fields[1] == "PARTIAL" && fields[2] == "RESULTS" && fields[3] == "***" {
		f |= obs.Partial
	}
	return f, nil
}

func stripVia(fields []string) []string {
	for i, f := range fields {
		if f == "via" {
			return fields[:i]
		}
	}
	return fields
}
