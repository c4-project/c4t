// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"errors"
	"fmt"

	"github.com/MattWindsor91/c4t/internal/subject/compilation"
)

var (
	// ErrNoBin occurs when a successful compile result	has no binary path attached.
	ErrNoBin = errors.New("no binary in compile result")

	// ErrParserNil occurs when a runner config doesn't specify an observation parser.
	ErrParserNil = errors.New("obs-parser nil")
)

// Error is the main error type returned by the runner.
// It wraps an inner error with various pieces of context.
type Error struct {
	// Stage is a representation of the part of the runner that went wrong.
	Stage string

	// Compilation is the name of the compilation whose run caused the error.
	Compilation compilation.Name

	// Inner is the inner error, if any, that caused this error.
	Inner error
}

// Error implements the error protocol for Error.
func (e Error) Error() string {
	return fmt.Sprintf("while %s compilation %s: %s",
		e.Stage, e.Compilation, e.Inner.Error(),
	)
}

// Unwrap unwraps an Error, returning its inner error.
func (e Error) Unwrap() error {
	return e.Inner
}
