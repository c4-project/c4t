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

	// Num is the number of status flags.
	Num
	// FirstBad refers to the first status that is neither OK nor 'unknown'.
	FirstBad = Flagged
)

var (
	// ErrBad occurs when FromString encounters an unknown status string.
	ErrBad = errors.New("bad status")

	// Strings enumerates string equivalents for each Status.
	Strings = [Num]string{
		"unknown",
		"ok",
		"flagged",
		"compile/fail",
		"compile/timeout",
		"run/fail",
		"run/timeout",
	}
)

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
	for i, sc := range Strings {
		if strings.EqualFold(s, sc) {
			return Status(i), nil
		}
	}
	return Unknown, fmt.Errorf("%w: %q", ErrBad, s)
}

// String gets the string representation of a Status.
func (s Status) String() string {
	if len(Strings) <= int(s) || s < 0 {
		return "(BAD STATUS)"
	}
	return Strings[s]
}

// IsOk is true if, and only if, this status is StatusOk.
func (s Status) IsOk() bool {
	return s == Ok
}
