// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package singleobs contains observer implementations for use in the 'single-shot' c4t commands.
package singleobs

import (
	"log"

	"github.com/c4-project/c4t/internal/coverage"

	"github.com/c4-project/c4t/internal/director"

	"github.com/c4-project/c4t/internal/copier"
	"github.com/c4-project/c4t/internal/stage/mach/observer"

	"github.com/c4-project/c4t/internal/stage/perturber"

	"github.com/c4-project/c4t/internal/stage/planner"
	"github.com/c4-project/c4t/internal/subject/corpus/builder"
)

// DirectorInstance builds a list of director-instance compatible observers suitable for single-shot binaries.
//
// While it is very unlikely that these observers will be used in a director instance, since those have specific
// integrated observers, the director-instance interface covers almost every observer interface we want to implement
// for the single-shot binaries, and so we can derive the other observer constructors from it.
func DirectorInstance(l *log.Logger, verbose bool) []director.InstanceObserver {
	if !verbose {
		return nil
	}

	// The ordering is important here: we want log messages to appear _before_ progress bars.
	return []director.InstanceObserver{
		NewBar(),
		(*Logger)(l),
	}
}

// Planner builds a list of observers suitable for single-shot c4t planner binaries.
func Planner(l *log.Logger, verbose bool) []planner.Observer {
	// Director observers don't implement the planner, annoyingly.
	if !verbose {
		return nil
	}
	return []planner.Observer{NewBar(), (*Logger)(l)}
}

// Perturber builds a list of observers suitable for single-shot c4t perturber binaries.
func Perturber(l *log.Logger, verbose bool) []perturber.Observer {
	return director.LowerToPerturber(DirectorInstance(l, verbose))
}

// Builder builds a list of observers suitable for single-shot c4t corpus-builder binaries.
func Builder(l *log.Logger, verbose bool) []builder.Observer {
	return director.LowerToBuilder(DirectorInstance(l, verbose))
}

// Copier builds a list of observers suitable for observing file copies in single-shot binaries.
func Copier(l *log.Logger, verbose bool) []copier.Observer {
	return director.LowerToCopy(DirectorInstance(l, verbose))
}

// Mach builds a list of observers suitable for observing machine node actions in single-shot binaries.
func MachNode(l *log.Logger, verbose bool) []observer.Observer {
	return director.LowerToMach(DirectorInstance(l, verbose))
}

// Coverage builds a list of observers suitable for observing coverage testbed generation.
func Coverage(l *log.Logger, verbose bool) []coverage.Observer {
	// TODO(@MattWindsor91): add bar support?
	if !verbose {
		return []coverage.Observer{}
	}
	return []coverage.Observer{(*Logger)(l)}
}
