// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package herd

import (
	"fmt"
	"strconv"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/herdtools/parser"
)

// ParseStateCount parses a Herd state count.
func (h Herd) ParseStateCount(fields []string) (uint64, error) {
	if nf := len(fields); nf != 2 {
		return 0, fmt.Errorf("%w: expected two fields, got %d", parser.ErrBadStateCount, nf)
	}
	if f := fields[0]; f != "States" {
		return 0, fmt.Errorf("%w: expected first word to be 'State', got %q", parser.ErrBadStateCount, f)
	}
	return strconv.ParseUint(fields[1], 10, 64)
}

// ParseStateLine 'parses' a Herd state line.
// Herd state lines need no actual processing, and just get passed through verbatim.
func (h Herd) ParseStateLine(_ parser.TestType, fields []string) (*parser.StateLine, error) {
	return &parser.StateLine{Rest: fields}, nil
}
