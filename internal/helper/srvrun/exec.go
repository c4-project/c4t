// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package srvrun

import (
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/MattWindsor91/c4t/internal/model/service"
)

// ExecRunner represents runs of services using exec.
type ExecRunner struct {
	errw  io.Writer
	outw  io.Writer
	grace time.Duration
}

// NewExecRunner constructs an exec-based runner using the options in os.
func NewExecRunner(os ...ExecOption) *ExecRunner {
	e := new(ExecRunner)
	ExecOptions(os...)(e)
	return e
}

// WithStderr returns a new ExecRunner with standard error rerouted to w.
func (e ExecRunner) WithStderr(w io.Writer) service.Runner {
	StderrTo(w)(&e)
	return e
}

// WithGrace returns a new ExecRunner with the timeout grace period set as d.
func (e ExecRunner) WithGrace(d time.Duration) service.Runner {
	WithGrace(d)(&e)
	return e
}

// WithStdout returns a new ExecRunner with standard output rerouted to w.
func (e ExecRunner) WithStdout(w io.Writer) service.Runner {
	e.outw = w
	return e
}

func (e ExecRunner) hasGrace() bool {
	return runtime.GOOS != "windows" && 0 < e.grace
}

// Run runs the command specified by r using exec, on context ctx.
func (e ExecRunner) Run(ctx context.Context, r service.RunInfo) error {
	if e.hasGrace() {
		return e.runWithGrace(ctx, r)
	}
	c := exec.CommandContext(ctx, r.Cmd, r.Args...)
	c.Stderr = e.errw
	c.Stdout = e.outw
	return c.Run()
}

func (e ExecRunner) runWithGrace(ctx context.Context, r service.RunInfo) error {
	cdone := make(chan struct{})
	defer close(cdone)

	c := exec.Command(r.Cmd, r.Args...)
	c.Stderr = e.errw
	c.Stdout = e.outw
	if err := c.Start(); err != nil {
		return err
	}

	// https://github.com/golang/go/issues/22757#issuecomment-345009730
	go func() {
		select {
		case <-cdone:
			return
		case <-ctx.Done():
		}

		t := time.AfterFunc(e.grace, func() {
			_ = c.Process.Kill()
		})
		_ = c.Process.Signal(os.Interrupt)

		<-cdone
		t.Stop()
	}()

	return c.Wait()
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

// WithGrace sets a timeout grace period of d.
//
// If an ExecRunner's context closes and it has a timeout grace period, and the OS supports it, it will SIGTERM the
// program, wait d, and then SIGKILL.
func WithGrace(d time.Duration) ExecOption {
	return func(l *ExecRunner) {
		l.grace = d
	}
}
