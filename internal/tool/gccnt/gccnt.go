// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package gccnt implements the "GCCn't" wrapper over GCC-style compilers.
package gccnt

import (
	"context"
	"errors"
	"io"
	"sort"

	"github.com/c4-project/c4t/internal/serviceimpl/compiler/gcc"
)

var (
	// ErrAskedTo occurs when gccn't is asked to simulate a compiler failure.
	ErrAskedTo = errors.New("scheduled error")
	// ErrNoBin occurs when a gccn't runs without a binary set.
	ErrNoBin = errors.New("no gcc binary provided")
)

// Gccnt holds all configuration needed to perform a GCCn't run.
type Gccnt struct {
	// Conds stores the conditions for which gccn't will fail.
	Conds ConditionSet

	// Bin is the name of the 'real' compiler binary to run.
	Bin string

	// In gives the list of filepaths to pass to GCC.
	In []string

	// Out gives the filepath to which GCC is going to output.
	Out string

	// OptLevel contains the raw optimisation level string to run GCC with.
	OptLevel string

	// Mutant is, if nonzero, the current mutant.
	Mutant uint64

	// March specifies whether to pass -march to gcc, and, if so, what value.
	March string

	// Mcpu specifies whether to pass -mcpu to gcc, and, if so, what value.
	Mcpu string

	// Pthread specifies whether to pass -pthread to gcc.
	Pthread bool

	// Std specifies the standard to pass to gcc.
	Std string
}

// runner is the interface of low-level drivers that gccn't can use.
type runner interface {
	// Diverge should [pretend to] diverge, ie spin until and unless ctx cancels.
	Diverge(ctx context.Context) error

	// DoError should [pretend to] return an error that causes gccn't to fail.
	DoError() error

	// MutantSelect should log a mutant selection.
	MutantSelect(mutant uint64) error

	// MutantHit should log a mutant hit.
	MutantHit(mutant uint64) error

	// Init gives the Runner a chance to dump information about its configuration.
	Init(conds ConditionSet) error

	// RunGCC should run, or pretend to run, GCC with the given command bin and arguments args.
	RunGCC(ctx context.Context, bin string, args ...string) error
}

// Run runs gccn't on the given context, printing stdout onto outw and stderr onto errw.
func (g *Gccnt) Run(ctx context.Context, outw, errw io.Writer) error {
	return g.runOnRunner(ctx, &wetRunner{outw: outw, errw: errw})
}

func (g *Gccnt) runOnRunner(ctx context.Context, r runner) error {
	var err error

	if err = g.check(); err != nil {
		return err
	}

	g.Conds.sort()

	if err = g.handleMutant(r); err != nil {
		return err
	}

	if err := r.Init(g.Conds); err != nil {
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

func (g *Gccnt) handleMutant(r runner) error {
	if g.Mutant == 0 {
		return nil
	}
	var err error
	if err = r.MutantSelect(g.Mutant); err != nil {
		return err
	}
	if !g.mutantHit() {
		return nil
	}
	return r.MutantHit(g.Mutant)
}

func (g *Gccnt) mutantHit() bool {
	c := g.Conds
	for _, p := range []uint64{c.MutHitPeriod, c.Diverge.MutPeriod, c.Error.MutPeriod} {
		if g.mutantActive(p) {
			return true
		}
	}
	return false
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
	args = gcc.AddStringArg(args, "std=", g.Std)
	args = gcc.AddStringArg(args, "march=", g.March)
	args = gcc.AddStringArg(args, "mcpu=", g.Mcpu)

	if g.Pthread {
		args = append(args, "-pthread")
	}
	return append(args, g.In...)
}

// shouldDiverge checks whether gccn't should diverge.
func (g *Gccnt) shouldDiverge() bool {
	return g.should(g.Conds.Diverge)
}

// shouldError checks whether gccn't should run an error.
func (g *Gccnt) shouldError() bool {
	return g.should(g.Conds.Error)
}

func (g *Gccnt) should(c Condition) bool {
	return g.mutantActive(c.MutPeriod) || present(g.OptLevel, c.Opts)
}

func (g *Gccnt) mutantActive(period uint64) bool {
	return 0 < g.Mutant && 0 < period && g.Mutant%period == 0
}

func present(x string, xs []string) bool {
	i := sort.SearchStrings(xs, x)
	return i < len(xs) && xs[i] == x
}
