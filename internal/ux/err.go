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

	backend2 "github.com/c4-project/c4t/internal/model/service/backend"
	"github.com/c4-project/c4t/internal/stage/lifter"
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

	handleTopExitError(err)
	handleTopLiftingError(err)

	errln(err)
	os.Exit(1)
}

func handleTopExitError(err error) {
	var perr *exec.ExitError
	if errors.As(err, &perr) {
		logTopExitError(perr)
		os.Exit(perr.ExitCode())
	}
}

func logTopExitError(perr *exec.ExitError) {
	errln("child process encountered error:")
	errln(perr)
	if len(perr.Stderr) == 0 {
		return
	}
	errln("(any captured stderr follows)")
	errln(string(perr.Stderr))
}

func handleTopLiftingError(err error) {
	var lerr *lifter.Error
	if errors.As(err, &lerr) {
		logTopLiftingError(lerr)
		os.Exit(1)
	}
}

func logTopLiftingError(lerr *lifter.Error) {
	errln("Couldn't lift subject", lerr.SubjectName(), "with backend", lerr.ServiceName)
	errln(" - target architecture:", lerr.Job.Arch)
	errln(" - input type:", lerr.Job.In.Source)
	errln(" - output type:", lerr.Job.Out.Target)
	errln(lerr.Inner)

	if errors.Is(lerr.Inner, backend2.ErrNotSupported) {
		errln()
		errln("This backend doesn't support this particular type of lift operation.")
		errln("Check configuration/arguments - maybe the wrong backend has been selected?")
	}
}

func errln(a ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
}
