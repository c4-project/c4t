// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package rmach handles invoking the 'mach' controller in a network-transparent manner with the act-tester-mach binary.
package rmach

import (
	"github.com/1set/gut/ystring"
	"github.com/MattWindsor91/act-tester/internal/controller/rmach/runner"
)

// Invoker runs the machine-runner, through SSH if needed.
type Invoker struct {
	// dirLocal is the filepath to the directory to which local outcomes from this rmach run will appear.
	dirLocal string
	// invoker tells the remote-machine controller which arguments to send to the machine binary.
	invoker runner.InvocationGetter
	// observers is the set of observers listening for file copying and remote corpus manipulations.
	observers ObserverSet
	// rfac governs how the invoker will run the machine node when given a plan to invoke.
	rfac runner.Factory
}

// New constructs a new Mach with ssh configuration ssh (if any) and local directory ldir.
func New(ldir string, inv runner.InvocationGetter, o ...Option) (*Invoker, error) {
	if err := check(ldir, inv); err != nil {
		return nil, err
	}

	invoker := Invoker{dirLocal: ldir, invoker: inv, rfac: runner.LocalFactory(ldir)}
	if err := Options(o...)(&invoker); err != nil {
		return nil, err
	}
	return &invoker, nil
}

func check(ldir string, inv runner.InvocationGetter) error {
	if ystring.IsBlank(ldir) {
		return ErrDirEmpty
	}
	if inv == nil {
		return ErrInvokerNil
	}
	return nil
}
