// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package herdtools

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/job"
	"github.com/MattWindsor91/act-tester/internal/pkg/model/service"
)

// Herd describes the parts of a Backend invocation that are specific to Herd.
type Herd struct{}

var ErrNotSupported = errors.New("service doesn't support action")

// ParseStateCount parses a Herd state count.
func (h Herd) ParseStateCount(fields []string) (uint64, error) {
	if nf := len(fields); nf != 2 {
		return 0, fmt.Errorf("%w: expected two fields, got %d", ErrBadStateCount, nf)
	}
	if f := fields[0]; f != "States" {
		return 0, fmt.Errorf("%w: expected first word to be 'State', got %q", ErrBadStateCount, f)
	}
	return strconv.ParseUint(fields[1], 10, 64)
}

// ParseStateLine 'parses' a Herd state line.
// Herd state lines need no actual processing, and just get passed through verbatim.
func (h Herd) ParseStateLine(_ TestType, fields []string) (*StateLine, error) {
	return &StateLine{Rest: fields}, nil
}

// Args deduces the appropriate arguments for running Herd on job j, with the merged run information r.
func (h Herd) Args(_ job.Harness, _ service.RunInfo) ([]string, error) {
	// TODO(@MattWindsor91): once we extend this to deal with non-harness jobs, add functionality here.
	return nil, fmt.Errorf("%w: harness making", ErrNotSupported)
}
