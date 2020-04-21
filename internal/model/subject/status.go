// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/model/obs"
)

// Status is the type of completed-run statuses.
type Status int

const (
	// StatusUnknown represents an unknown status.
	StatusUnknown Status = iota
	// StatusOk indicates that a run completed successfully without incident.
	StatusOk
	// StatusFlagged indicates that a run completed successfully, but its observation was interesting.
	// Usually this means a counter-example occurred.
	StatusFlagged
	// StatusCompileFail indicates that a run failed because of the compilation failing.
	StatusCompileFail
	// StatusCompileTimeout indicates that a run failed because the compilation timed out.
	StatusCompileTimeout // TODO(@MattWindsor91): use
	// StatusCompileFail indicates that a run failed directly.
	StatusRunFail // TODO(@MattWindsor91): use
	// StatusRunTimeout indicates that a run timed out.
	StatusRunTimeout

	// NumStatus is the number of status flags.
	NumStatus
	// FirstBadStatus refers to the first status that is neither OK nor 'unknown'.
	FirstBadStatus = StatusFlagged
)

var (
	// ErrBadStatus occurs when StatusOfString encounters an unknown status string.
	ErrBadStatus = errors.New("bad status")

	// StatusStrings enumerates string equivalents for each Status.
	StatusStrings = [NumStatus]string{
		"unknown",
		"ok",
		"flagged",
		"compile/fail",
		"compile/timeout",
		"run/fail",
		"run/timeout",
	}
)

// StatusOfCompileError tries to see if err represents a non-fatal issue such as a timeout or process error.
// If so, it converts that error to a status and returns it alongside nil.
// Otherwise, it propagates the error forwards.
func StatusOfCompileError(err error) (Status, error) {
	return statusOfError(err, StatusCompileTimeout, StatusCompileFail)
}

// StatusOfRunError tries to see if err represents a non-fatal issue such as a timeout or process error.
// If so, it converts that error to a status and returns it alongside nil.
// Otherwise, it propagates the error forwards.
func StatusOfRunError(err error) (Status, error) {
	return statusOfError(err, StatusRunTimeout, StatusRunFail)
}

func statusOfError(err error, timeout, fail Status) (Status, error) {
	var ee *exec.ExitError
	switch {
	case err == nil:
		return StatusOk, nil
	case errors.Is(err, context.DeadlineExceeded):
		return timeout, nil
	case errors.As(err, &ee):
		return fail, nil
	default:
		return StatusUnknown, err
	}
}

// StatusOfString tries to resolve s to a status code.
func StatusOfString(s string) (Status, error) {
	for i, sc := range StatusStrings {
		if strings.EqualFold(s, sc) {
			return Status(i), nil
		}
	}
	return StatusUnknown, fmt.Errorf("%w: %q", ErrBadStatus, s)
}

// String gets the string representation of a Status.
func (s Status) String() string {
	if len(StatusStrings) <= int(s) || s < 0 {
		return "(BAD STATUS)"
	}
	return StatusStrings[s]
}

// MarshalText marshals a Status to text via its string representation.
func (s Status) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText unmarshals a Status from text via its string representation.
func (s *Status) UnmarshalText(text []byte) error {
	var err error
	*s, err = StatusOfString(string(text))
	return err
}

// IsOk is true if, and only if, this status is StatusOk.
func (s Status) IsOk() bool {
	return s == StatusOk
}

// StatusOfObs determines the status of an observation o given various items of context.
// The error runErr should contain any error that occurred when running the binary giving the observation.
// StatusOfObs returns any error passed to it that it deems too fatal to represent in the status code.
func StatusOfObs(o *obs.Obs, runErr error) (Status, error) {
	if runErr != nil {
		return StatusOfRunError(runErr)
	}

	// TODO(@MattWindsor91): allow interestingness criteria
	if o.Unsat() {
		return StatusFlagged, nil
	}

	return StatusOk, nil
}
