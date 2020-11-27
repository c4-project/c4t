// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gccnt

import (
	"context"
	"fmt"
	"io"
)

// DryRun works out what GCCn't is going to do, then prints it onto errw.
func (g *Gccnt) DryRun(ctx context.Context, errw io.Writer) error {
	return g.runOnRunner(ctx, &dryRunner{w: errw})
}

// dryRunner wraps a Writer with some state on any error that has been caused by dumping.
type dryRunner struct {
	w   io.Writer
	err error
}

// Diverge pretends to diverge, but instead just says that it normally would.
func (d *dryRunner) Diverge(context.Context) error {
	d.dumpln("gccn't would diverge here")
	return d.err
}

// DoError pretends to error out, but instead just says that it normally would.
func (d *dryRunner) DoError() error {
	d.dumpln("gccn't would error here")
	return d.err
}

// Init logs the error and divergence optimisation levels.
func (d *dryRunner) Init(errorOpts []string, divergeOpts []string) error {
	for _, c := range []struct {
		name string
		opts []string
	}{
		{name: "an error", opts: errorOpts},
		{name: "divergence", opts: divergeOpts},
	} {
		if len(c.opts) == 0 {
			continue
		}

		d.dumpf("The following optimisation levels will trigger %s:", c.name)
		for _, o := range c.opts {
			d.dumpf(" %s", o)
		}
		d.dumpln()
	}
	return d.err
}

func (d *dryRunner) RunGCC(_ context.Context, bin string, args ...string) error {
	d.dumpf("invocation: %s", bin)
	for _, a := range args {
		d.dumpf(" %s", a)
	}
	d.dumpln()
	return d.err
}

// dumpf Fprintf-s to the dumper's writer if no error has yet occurred.
func (d *dryRunner) dumpf(format string, a ...interface{}) {
	if d.err != nil {
		return
	}
	_, d.err = fmt.Fprintf(d.w, format, a...)
}

// dumpln Fprintln-s to the dumper's writer if no error has yet occurred.
func (d *dryRunner) dumpln(a ...interface{}) {
	if d.err != nil {
		return
	}
	_, d.err = fmt.Fprintln(d.w, a...)
}
