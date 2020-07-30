// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package singleobs contains observer implementations for use in the 'single-shot' act-tester commands.
package singleobs

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/stage/perturber"

	"github.com/MattWindsor91/act-tester/internal/stage/invoker"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/stage/planner"
)

// Planner builds a list of observers suitable for single-shot act-tester planner binaries.
func Planner(l *log.Logger) []planner.Observer {
	// The ordering is important here: we want log messages to appear _before_ progress bars.
	return []planner.Observer{
		NewBar(),
		(*Logger)(l),
	}
}

// Perturber builds a list of observers suitable for single-shot act-tester planner binaries.
func Perturber(l *log.Logger) []perturber.Observer {
	// The ordering is important here: we want log messages to appear _before_ progress bars.
	return []perturber.Observer{
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
func Invoker(l *log.Logger) []invoker.Observer {
	// The ordering is important here: we want log messages to appear _before_ progress bars.
	return []invoker.Observer{
		NewBar(),
		(*Logger)(l),
	}
}
