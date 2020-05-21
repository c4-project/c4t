// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package litmus implements a wrapper over the herdtools7 'litmus7' program.
// This wrapper deals with various corner cases.

package litmus

import (
	"context"
	"errors"
	"io"
	"os/exec"

	"github.com/MattWindsor91/act-tester/internal/model"
)

var (
	// ErrConfigNil occurs when Run tries to run on a nil config.
	ErrConfigNil = errors.New("config nil")

	// ErrStatNil occurs when the config has a nil statistics dumper.
	ErrStatNil = errors.New("config stat dumper nil")

	// ErrNoCArch occurs when the output directory is empty.
	ErrNoCArch = errors.New("need carch")
)

// Litmus is the configuration required to run the litmus shim.
type Litmus struct {
	// Stat extracts statistics from litmus files.
	// These statistics then switch on various fixes.
	Stat model.StatDumper

	// Err is the writer to which stderr output should be written.
	Err io.Writer

	// CArch is the architecture that the litmus shim should target.
	// It corresponds to Litmus's 'carch' argument.
	CArch string

	// Fixset is the set of enabled fixes.
	// It is part of the config to allow the forcing of fixes that the shim would otherwise deem unnecessary.
	Fixset Fixset

	// Verbose toggles various 'verbose' dumping actions.
	Verbose bool

	// Pathset is the set of specified paths for this litmus invocation.
	Pathset Pathset
}

// Run runs the litmus wrapper according to the configuration c.
func (l *Litmus) Run(ctx context.Context) error {
	if err := l.check(); err != nil {
		return err
	}

	if err := l.probeFixes(ctx); err != nil {
		return err
	}

	if err := l.runLitmus(); err != nil {
		return err
	}

	return l.patch()
}

// check checks that the configuration makes sense.
func (l *Litmus) check() error {
	if l == nil {
		return ErrConfigNil
	}
	if l.Stat == nil {
		return ErrStatNil
	}
	if l.CArch == "" {
		return ErrNoCArch
	}
	return l.Pathset.Check()
}

// probeFixes checks to see if there are any fixes needed for the input.
func (l *Litmus) probeFixes(ctx context.Context) error {
	var s model.Statset
	if err := l.Stat.DumpStats(ctx, &s, l.Pathset.FileIn); err != nil {
		return err
	}
	l.Fixset.PopulateFromStats(&s)

	if l.Verbose {
		return l.Fixset.Dump(l.Err)
	}
	return nil
}

// runLitmus actually runs Litmus.
func (l *Litmus) runLitmus() error {
	cmd := exec.Command("litmus7", l.litmusArgs()...)

	cmd.Stderr = l.Err
	return cmd.Run()
}

// litmusArgs works out the argument vector for Litmus.
func (l *Litmus) litmusArgs() []string {
	args := l.Fixset.Args()
	args = append(args, "-carch", l.CArch, "-c11", "true")
	args = append(args, l.Pathset.Args()...)
	return args
}
