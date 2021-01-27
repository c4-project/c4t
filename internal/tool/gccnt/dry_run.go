// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gccnt

import (
	"context"
	"fmt"
	"io"

	"github.com/c4-project/c4t/internal/mutation"
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
func (d *dryRunner) Init(conds ConditionSet) error {
	for _, c := range []struct {
		name string
		cnd  Condition
	}{
		// Not a map, to make sure we get a consistent ordering.
		{name: "divergence", cnd: conds.Diverge},
		{name: "an error", cnd: conds.Error},
	} {
		d.dumpCondition(c.name, c.cnd)
	}
	return d.err
}

func (d *dryRunner) dumpCondition(name string, c Condition) {
	d.dumpMut(name, c.MutPeriod)
	d.dumpOpts(name, c.Opts)
}

func (d *dryRunner) dumpMut(name string, period uint64) {
	if period == 0 {
		return
	}
	d.dumpf("Mutation numbers that are multiples of %d will trigger %s\n", period, name)
}

func (d *dryRunner) dumpOpts(name string, opts []string) {
	if len(opts) == 0 {
		return
	}
	d.dumpf("The following optimisation levels will trigger %s:", name)
	for _, o := range opts {
		d.dumpf(" %s", o)
	}
	d.dumpln()
}

// MutantHit dumps the mutation hit stanza.
func (d *dryRunner) MutantHit(n uint64) error {
	d.dumpln(mutation.MutantHitPrefix, n)
	return d.err
}

// MutantSelect dumps the mutation selection stanza.
func (d *dryRunner) MutantSelect(n uint64) error {
	d.dumpln(mutation.MutantSelectPrefix, n)
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
