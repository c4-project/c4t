// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package ux

import (
	"errors"
	"os/exec"

	"github.com/sirupsen/logrus"
)

// LogTopError logs err if non-nil, in a 'fatal top-level error' sort of way.
func LogTopError(err error) {
	if err == nil {
		return
	}

	var perr *exec.ExitError
	if errors.As(err, &perr) {
		logrus.WithError(err).WithField("stderr", string(perr.Stderr)).Errorln("a child process encountered an error")
		return
	}

	logrus.WithError(err).Errorln("fatal error")
}
