// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package gccnt implements the "GCCn't" wrapper over GCC-style compilers.
package gccnt

import (
	"context"
	"errors"
	"io"
	"sort"
)

var (
	// ErrAskedTo occurs when gccn't is asked to simulate a compiler failure.
	ErrAskedTo = errors.New("scheduled error")
	// ErrNoBin occurs when a gccn't runs without a binary set.
	ErrNoBin = errors.New("no gcc binary provided")
)

// Gccnt holds all configuration needed to perform a GCCn't run.
type Gccnt struct {
	// Bin is the name of the 'real' compiler binary to run.
	Bin string

	// In gives the list of filepaths to pass to GCC.
	In []string

	// Out gives the filepath to which GCC is going to output.
	Out string

	// OptLevel contains the raw optimisation level string to run GCC with.
	OptLevel string

	// ErrorOpts contains the optimisation levels at which gccn't will error.
	ErrorOpts []string

	// DivergeOpts contains the optimisation levels at which gccn't will diverge.
	DivergeOpts []string

	// Std specifies the standard to pass to gcc.
	Std string

	// Pthread specifies whether to pass -pthread to gcc.
	Pthread bool
}

// runner is the interface of low-level drivers that gccn't can use.
type runner interface {
	// Diverge should [pretend to] diverge, ie spin until and unless ctx cancels.
	Diverge(ctx context.Context) error

	// DoError should [pretend to] return an error that causes gccn't to fail.
	DoError() error

	// Init gives the Runner a chance to dump information about its configuration.
	Init(errorOpts []string, divergeOpts []string) error

	// RunGCC should run, or pretend to run, GCC with the given command bin and arguments args.
	RunGCC(ctx context.Context, bin string, args ...string) error
}

// Run runs gccn't on the given context, printing stdout onto outw and stderr onto errw.
func (g *Gccnt) Run(ctx context.Context, outw, errw io.Writer) error {
	return g.runOnRunner(ctx, &wetRunner{outw: outw, errw: errw})
}

func (g *Gccnt) runOnRunner(ctx context.Context, r runner) error {
	if err := g.check(); err != nil {
		return err
	}

	sort.Strings(g.DivergeOpts)
	sort.Strings(g.ErrorOpts)

	if err := r.Init(g.ErrorOpts, g.DivergeOpts); err != nil {
		return err
	}

	switch {
	case g.shouldDiverge():
		return r.Diverge(ctx)
	case g.shouldError():
		return r.DoError()
	default:
		return r.RunGCC(ctx, g.Bin, g.args()...)
	}
}

func (g *Gccnt) check() error {
	if g.Bin == "" {
		return ErrNoBin
	}
	return nil
}

func (g *Gccnt) args() []string {
	args := []string{
		"-o", g.Out,
		"-O" + g.OptLevel,
	}
	if g.Std != "" {
		args = append(args, "-std="+g.Std)
	}
	if g.Pthread {
		args = append(args, "-pthread")
	}
	return append(args, g.In...)
}

// shouldDiverge checks whether gccn't should diverge.
func (g *Gccnt) shouldDiverge() bool {
	return g.should(g.DivergeOpts)
}

// shouldError checks whether gccn't should run an error.
func (g *Gccnt) shouldError() bool {
	return g.should(g.ErrorOpts)
}

func (g *Gccnt) should(opts []string) bool {
	i := sort.SearchStrings(opts, g.OptLevel)
	return i < len(opts) && opts[i] == g.OptLevel
}
