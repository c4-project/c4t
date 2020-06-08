// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package srvrun

import (
	"context"
	"io"
	"os/exec"

	"github.com/MattWindsor91/act-tester/internal/model/service"
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

/*
// StdoutTo redirects standard output of any commands run by the runner to w.
func StdoutTo(w io.Writer) ExecOption {
	return func(l *ExecRunner) {
		l.outw = w
	}
}
*/
