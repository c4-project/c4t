// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package invoker handles invoking the 'mach' stage in a network-transparent manner with the act-tester-mach binary.
package invoker

import (
	"github.com/1set/gut/ystring"
	"github.com/MattWindsor91/act-tester/internal/copier"
	"github.com/MattWindsor91/act-tester/internal/quantity"
	"github.com/MattWindsor91/act-tester/internal/stage/invoker/runner"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer"
)

// Invoker runs the machine-runner, through SSH if needed.
type Invoker struct {
	// ldir is the local directory to which machine node files are to be copied.
	ldir string
	// baseQuantities contains the base quantity set, to be overridden by any other quantities given.
	baseQuantities quantity.MachNodeSet
	// copyObservers is the set of observers listening for file copying.
	copyObservers []copier.Observer
	// machObservers is the set of observers listening for remote corpus manipulations.
	machObservers []observer.Observer
	// rfac governs how the invoker will run the machine node when given a plan to invoke.
	rfac runner.Factory
}

// New constructs a new Invoker with local directory ldir, runner factory fac, and options o.
func New(ldir string, fac runner.Factory, o ...Option) (*Invoker, error) {
	if ystring.IsBlank(ldir) {
		return nil, ErrDirEmpty
	}

	invoker := Invoker{ldir: ldir, rfac: fac}
	if err := Options(o...)(&invoker); err != nil {
		return nil, err
	}
	return &invoker, nil
}
