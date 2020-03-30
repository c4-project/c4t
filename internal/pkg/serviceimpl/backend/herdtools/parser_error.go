// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package herdtools

import "errors"

var (
	// ErrNoImpl occurs if the parser's impl is nil.
	ErrNoImpl = errors.New("parser not supplied with valid implementation")
	// ErrBadState occurs if the parser somehow gets into an unknown state.
	ErrBadState = errors.New("unknown state")
	// ErrBadTransition occurs if the parser somehow performs a bad transition.
	ErrBadTransition = errors.New("invalid state transition")

	// ErrBadStateCount occurs if the number of states is invalid.
	ErrBadStateCount = errors.New("bad state count")
	// ErrInputEmpty occurs if the input to the parser runs out before anything is parsed.
	ErrInputEmpty = errors.New("input was empty")
	// ErrNoTest occurs if the input to the parser runs out before the start of a test is parsed.
	ErrNoTest = errors.New("input ended before reaching test")
	// ErrNoStates occurs if the input to the parser runs out before any states are parsed.
	ErrNoStates = errors.New("input ended with no state block")
	// ErrNotEnoughStates occurs if the input to the parser runs out midway through the expected number of states.
	ErrNotEnoughStates = errors.New("input ended while expecting more states")
	// ErrNoSummary occurs if the input to the parser runs out before the summary is parsed.
	ErrNoSummary = errors.New("input ended while expecting summary")

	// ErrBadTestType occurs if the test-type line is badly formed.
	ErrBadTestType = errors.New("malformed test-type line")
	// ErrBadSummary occurs if the summary line is badly formed.
	ErrBadSummary = errors.New("malformed summary line")
	// ErrBadStateLine occurs if the state line is badly formed.
	// It may be used by BackendImpl implementations.
	ErrBadStateLine = errors.New("malformed state line")
)
