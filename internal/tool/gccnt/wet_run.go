// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gccnt

import (
	"context"
	"io"
	"os/exec"
)

// wetRunner implements the low-level running logic of gccn't for real, as opposed to dry-running.
type wetRunner struct {
	outw, errw io.Writer
}

// Diverge spins until and unless ctx cancels.
func (f *wetRunner) Diverge(ctx context.Context) error {
	for range ctx.Done() {
	}
	return ctx.Err()
}

// DoError returns an error that causes gccn't to fail.
func (f *wetRunner) DoError() error {
	return ErrAskedTo
}

// Init does nothing.
func (f *wetRunner) Init(_, _ []string) error {
	return nil
}

// RunGCC runs GCC.
func (f *wetRunner) RunGCC(ctx context.Context, bin string, args ...string) error {
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stderr = f.errw
	cmd.Stdout = f.outw
	return cmd.Run()
}
