// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package status

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Status is the type of completed-run statuses.
type Status int

const (
	// Unknown represents an unknown status.
	Unknown Status = iota
	// Ok indicates that a run completed successfully without incident.
	Ok
	// Filtered indicates that a run would have failed, but it has been caught by a filter.
	Filtered
	// Flagged indicates that a run completed successfully, but its observation was interesting.
	// Usually this means a counter-example occurred.
	Flagged
	// CompileFail indicates that a run failed because of the compilation failing.
	CompileFail
	// CompileTimeout indicates that a run failed because the compilation timed out.
	CompileTimeout
	// RunFail indicates that a run failed directly.
	RunFail
	// RunTimeout indicates that a run timed out.
	RunTimeout

	// FirstBad refers to the first status that represents an unwanted outcome.
	FirstBad = Flagged
	// Last is the last valid status.
	Last = RunTimeout
)

//go:generate stringer -type=Status

// ErrBad occurs when FromString encounters an unknown status string.
var ErrBad = errors.New("bad status")

// FromCompileError tries to see if err represents a non-fatal issue such as a timeout or process error.
// If so, it converts that error to a status and returns it alongside nil.
// Otherwise, it propagates the error forwards.
func FromCompileError(err error) (Status, error) {
	return statusOfError(err, CompileTimeout, CompileFail)
}

// FromRunError tries to see if err represents a non-fatal issue such as a timeout or process error.
// If so, it converts that error to a status and returns it alongside nil.
// Otherwise, it propagates the error forwards.
func FromRunError(err error) (Status, error) {
	return statusOfError(err, RunTimeout, RunFail)
}

func statusOfError(err error, timeout, fail Status) (Status, error) {
	var ee *exec.ExitError
	switch {
	case err == nil:
		return Ok, nil
	case errors.Is(err, context.DeadlineExceeded):
		return timeout, nil
	case errors.As(err, &ee):
		return fail, nil
	default:
		return Unknown, err
	}
}

// FromString tries to resolve s to a status code.
func FromString(s string) (Status, error) {
	for i := Unknown; i <= Last; i++ {
		if strings.EqualFold(s, i.String()) {
			return Status(i), nil
		}
	}
	return Unknown, fmt.Errorf("%w: %q", ErrBad, s)
}

// IsOk is true if, and only if, this status is StatusOk.
func (i Status) IsOk() bool {
	return i == Ok
}

// IsInBounds is true if, and only if, this status is in bounds.
func (i Status) IsInBounds() bool {
	return Ok <= i && i <= Last
}

// IsBad is true if, and only if, this status is in-bounds, not OK, and not unknown.
func (i Status) IsBad() bool {
	return FirstBad <= i && i <= Last
}

// CountsForTiming is true if this status should be logged in compiler/run timing data.
//
// Any compile or run that is not filtered and executes to completion counts for timing purposes (so, ok or flagged
// compiles/runs).
func (i Status) CountsForTiming() bool {
	return !(FlagFail | FlagTimeout).MatchesStatus(i)
}
