// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package ux

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// LogTopError logs err if non-nil, in a 'fatal top-level error' sort of way.
func LogTopError(err error) {
	if err == nil {
		return
	}

	if errors.Is(err, context.Canceled) {
		// Assume that a top-level cancellation is user-specified.
		return
	}

	var perr *exec.ExitError
	if errors.As(err, &perr) {
		logTopExitError(err, perr)
		os.Exit(perr.ExitCode())
	}

	logTopNormalError(err)
	os.Exit(1)
}

func logTopExitError(err error, perr *exec.ExitError) {
	_, _ = fmt.Fprintln(os.Stderr, "child process encountered error:")
	_, _ = fmt.Fprintln(os.Stderr, err)
	if len(perr.Stderr) == 0 {
		return
	}
	_, _ = fmt.Fprintln(os.Stderr, "(any captured stderr follows)")
	_, _ = fmt.Fprintln(os.Stderr, string(perr.Stderr))
}

func logTopNormalError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
}
