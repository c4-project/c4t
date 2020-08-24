// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package singleobs contains observer implementations for use in the 'single-shot' act-tester commands.
package singleobs

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/copier"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer"

	"github.com/MattWindsor91/act-tester/internal/stage/perturber"

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

// Copier builds a list of observers suitable for observing file copies in single-shot binaries.
func Copier(l *log.Logger) []copier.Observer {
	// See above.
	return []copier.Observer{
		NewBar(),
		(*Logger)(l),
	}
}

// Mach builds a list of observers suitable for observing machine node actions in single-shot binaries.
func MachNode(l *log.Logger) []observer.Observer {
	// See above.
	return []observer.Observer{
		NewBar(),
		(*Logger)(l),
	}
}
