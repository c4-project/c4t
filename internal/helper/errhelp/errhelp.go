// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package errhelp contains error helper functions.
package errhelp

import (
	"context"
	"errors"
)

// FirstError gets the first error in errs that is non-nil, or nil if none exist.
func FirstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// TimeoutOrFirstError returns ctx.Err() if it is a deadline-exceeded error, or the first error in errs that is non-nil,
// if any exist.
func TimeoutOrFirstError(ctx context.Context, errs ...error) error {
	ctxErr := ctx.Err()
	if ctxErr != nil && errors.Is(ctxErr, context.DeadlineExceeded) {
		return ctxErr
	}
	return FirstError(errs...)
}
