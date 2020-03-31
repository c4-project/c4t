// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import "context"

// CheckDone checks, without blocking, to see if ctx has been cancelled as of this instant in time.
// If so, it propagates the context's error value; if not, it returns nil.
func CheckDone(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
