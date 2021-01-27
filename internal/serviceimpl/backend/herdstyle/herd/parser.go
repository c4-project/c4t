// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package herd

import (
	"fmt"
	"strconv"

	"github.com/c4-project/c4t/internal/subject/obs"

	"github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle/parser"
)

// ParseStateCount parses a Herd state count.
func (Herd) ParseStateCount(fields []string) (k uint64, ok bool, err error) {
	nf := len(fields)
	// This might not be a state count line.
	// In practice, Herd always follows the preamble with a state count, but Rmem, which otherwise follows the Herd
	// syntax here, does.
	if nf == 0 || fields[0] != "States" {
		return 0, false, nil
	}
	// At this point, we're expecting a state count line.
	if nf != 2 {
		return 0, false, fmt.Errorf("%w: expected two fields, got %d", parser.ErrBadStateCount, nf)
	}
	if f := fields[0]; f != "States" {
		return 0, false, fmt.Errorf("%w: expected first word to be 'State', got %q", parser.ErrBadStateCount, f)
	}
	k, err = strconv.ParseUint(fields[1], 10, 64)
	return k, true, err
}

// ParseStateLine 'parses' a Herd state line.
// Herd state lines need no actual processing, and just get passed through verbatim.
func (Herd) ParseStateLine(_ parser.TestType, fields []string) (*parser.StateLine, error) {
	return &parser.StateLine{Rest: fields}, nil
}

// ParsePreTestLine does nothing, as pre-Test lines have no meaning in Herd.
func (Herd) ParsePreTestLine([]string) (obs.Flag, error) {
	return 0, nil
}
