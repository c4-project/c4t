// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package gccnt implements the "GCCn't" wrapper over GCC-style compilers.
package gccnt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

var (
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
}

// DryRun works out what GCCn't is going to do, then prints it onto errw.
func (g *Gccnt) DryRun(errw io.Writer) error {
	if err := g.check(); err != nil {
		return err
	}

	d := dumper{w: errw}

	g.dumpInvocation(&d)

	return d.err
}

// dumper wraps a Writer with some state on any error that has been caused by dumping.
type dumper struct {
	w   io.Writer
	err error
}

func (g *Gccnt) dumpInvocation(d *dumper) {
	d.dumpf("invocation: %s", g.Bin)
	for _, a := range g.args() {
		d.dumpf(" %s", a)
	}
	d.dumpln()
}

// dumpf Fprintf-s to the dumper's writer if no error has yet occurred.
func (d *dumper) dumpf(format string, a ...interface{}) {
	if d.err != nil {
		return
	}
	_, d.err = fmt.Fprintf(d.w, format, a...)
}

// dumpln Fprintln-s to the dumper's writer if no error has yet occurred.
func (d *dumper) dumpln(a ...interface{}) {
	if d.err != nil {
		return
	}
	_, d.err = fmt.Fprintln(d.w, a...)
}

// Run runs gccn't on the given context, printing stdout onto outw and stderr onto errw.
func (g *Gccnt) Run(ctx context.Context, outw, errw io.Writer) error {
	if err := g.check(); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, g.Bin, g.args()...)
	cmd.Stderr = outw
	cmd.Stdout = errw
	return cmd.Run()
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
	return append(args, g.In...)
}
