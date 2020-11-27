// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package srvrun

import (
	"context"
	"io"
	"os/exec"

	"github.com/MattWindsor91/c4t/internal/model/service"
)

// ExecRunner represents runs of services using exec.
type ExecRunner struct {
	errw io.Writer
	outw io.Writer
}

// NewExecRunner constructs an exec-based runner using the options in os.
func NewExecRunner(os ...ExecOption) *ExecRunner {
	e := new(ExecRunner)
	ExecOptions(os...)(e)
	return e
}

// WithStderr returns a new ExecRunner with standard error rerouted to w.
func (e ExecRunner) WithStderr(w io.Writer) service.Runner {
	return ExecRunner{outw: e.outw, errw: w}
}

// WithStdout returns a new ExecRunner with standard output rerouted to w.
func (e ExecRunner) WithStdout(w io.Writer) service.Runner {
	return ExecRunner{outw: w, errw: e.errw}
}

// Run runs the command specified by r using exec, on context ctx.
func (e ExecRunner) Run(ctx context.Context, r service.RunInfo) error {
	c := exec.CommandContext(ctx, r.Cmd, r.Args...)
	c.Stderr = e.errw
	c.Stdout = e.outw
	return c.Run()
}

// ExecOption is the type of options used when constructing a ExecRunner.
type ExecOption func(*ExecRunner)

// ExecOptions applies each option in os.
func ExecOptions(os ...ExecOption) ExecOption {
	return func(l *ExecRunner) {
		for _, o := range os {
			o(l)
		}
	}
}

// StderrTo redirects standard error of any commands run by the runner to w.
func StderrTo(w io.Writer) ExecOption {
	return func(l *ExecRunner) {
		l.errw = w
	}
}
