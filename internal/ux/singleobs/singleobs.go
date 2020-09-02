// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package singleobs contains observer implementations for use in the 'single-shot' act-tester commands.
package singleobs

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/copier"
	dobserver "github.com/MattWindsor91/act-tester/internal/director/observer"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer"

	"github.com/MattWindsor91/act-tester/internal/stage/perturber"

	"github.com/MattWindsor91/act-tester/internal/stage/planner"
	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"
)

// DirectorInstance builds a list of director-instance compatible observers suitable for single-shot binaries.
//
// While it is very unlikely that these observers will be used in a director instance, since those have specific
// integrated observers, the director-instance interface covers every other observer interface we want to implement
// for the single-shot binaries, and so we can derive the other observer constructors from it.
func DirectorInstance(l *log.Logger) []dobserver.Instance {
	// The ordering is important here: we want log messages to appear _before_ progress bars.
	return []dobserver.Instance{
		NewBar(),
		(*Logger)(l),
	}
}

// Planner builds a list of observers suitable for single-shot act-tester planner binaries.
func Planner(l *log.Logger) []planner.Observer {
	return dobserver.LowerToPlanner(DirectorInstance(l))
}

// Perturber builds a list of observers suitable for single-shot act-tester planner binaries.
func Perturber(l *log.Logger) []perturber.Observer {
	return dobserver.LowerToPerturber(DirectorInstance(l))
}

// Builder builds a list of observers suitable for single-shot act-tester corpus-builder binaries.
func Builder(l *log.Logger) []builder.Observer {
	return dobserver.LowerToBuilder(DirectorInstance(l))
}

// Copier builds a list of observers suitable for observing file copies in single-shot binaries.
func Copier(l *log.Logger) []copier.Observer {
	return dobserver.LowerToCopy(DirectorInstance(l))
}

// Mach builds a list of observers suitable for observing machine node actions in single-shot binaries.
func MachNode(l *log.Logger) []observer.Observer {
	return dobserver.LowerToMach(DirectorInstance(l))
}
