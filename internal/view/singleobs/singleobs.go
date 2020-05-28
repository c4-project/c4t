// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package singleobs contains observer implementations for use in the 'single-shot' act-tester commands.
package singleobs

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/controller/rmach"

	"github.com/MattWindsor91/act-tester/internal/controller/planner"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
)

// Planner builds a list of observers suitable for single-shot act-tester planner binaries.
func Planner(l *log.Logger) []planner.Observer {
	// The ordering is important here: we want log messages to appear _before_ progress bars.
	return []planner.Observer{
		NewBar(),
		(*Logger)(l),
	}
}

// Builder builds a list of observers suitable for single-shot act-tester corpus-builder binaries.
func Builder(l *log.Logger) []builder.Observer {
	// See above.
	return []builder.Observer{
		NewBar(),
		(*Logger)(l),
	}
}

// Invoker builds a list of observers suitable for single-shot act-tester remote-mach binaries.
func RMach(l *log.Logger) []rmach.Observer {
	// The ordering is important here: we want log messages to appear _before_ progress bars.
	return []rmach.Observer{
		NewBar(),
		(*Logger)(l),
	}
}
