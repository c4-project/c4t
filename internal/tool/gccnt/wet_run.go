// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gccnt

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/c4-project/c4t/internal/mutation"
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
func (f *wetRunner) Init(_ ConditionSet) error {
	return nil
}

// MutantHit dumps the mutation hit stanza.
func (f *wetRunner) MutantHit(n uint64) error {
	_, err := fmt.Fprintln(f.errw, mutation.MutantHitPrefix, n)
	return err
}

// MutantSelect dumps the mutation selection stanza.
func (f *wetRunner) MutantSelect(n uint64) error {
	_, err := fmt.Fprintln(f.errw, mutation.MutantSelectPrefix, n)
	return err
}

// RunGCC runs GCC.
func (f *wetRunner) RunGCC(ctx context.Context, bin string, args ...string) error {
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stderr = f.errw
	cmd.Stdout = f.outw
	return cmd.Run()
}
