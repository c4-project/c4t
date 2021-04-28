// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package service

import (
	"bytes"
	"context"
	"io"
	"time"
)

// Runner is the interface of things that can run, or pretend to run, services.
type Runner interface {
	// WithStdout should return a new runner with the standard output overridden to w.
	WithStdout(w io.Writer) Runner

	// WithStderr should return a new runner with the standard error overridden to w.
	WithStderr(w io.Writer) Runner

	// WithGrace should return a new runner with the timeout grace period set to d.
	WithGrace(d time.Duration) Runner

	// Run runs r using context ctx.
	Run(ctx context.Context, r RunInfo) error
}

//go:generate mockery --name=Runner

// RunAndCaptureStdout runs r with ctx and ri, and captures its stdout to a string.
// This behaviour overrides any previous stdout capture for the runner.
func RunAndCaptureStdout(ctx context.Context, r Runner, ri RunInfo) (string, error) {
	var buf bytes.Buffer
	err := r.WithStdout(&buf).Run(ctx, ri)
	return buf.String(), err
}
