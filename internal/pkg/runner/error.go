// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// Error is the main error type returned by the runner.
// It wraps an inner error with various pieces of context.
type Error struct {
	// Stage is a representation of the part of the runner that went wrong.
	Stage string

	// Compiler is the ID of the compiler that produced the binary whose run caused the error.
	Compiler model.ID

	// Subject is the name of the subject that caused the error.
	Subject string

	// Inner is the inner error, if any, that caused this error.
	Inner error
}

// Error implements the error protocol for Error
func (e Error) Error() string {
	return fmt.Sprintf("while %s subject %s compile %s: %s",
		e.Stage, e.Subject, e.Compiler.String(), e.Inner.Error(),
	)
}

// Unwrap unwraps an Error, returning its inner error.
func (e Error) Unwrap() error {
	return e.Inner
}
