// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"github.com/MattWindsor91/c4t/internal/observing"
)

// Observer is the interface of types that can observe the progress of a coverage testbed generator.
type Observer interface {
	// OnCoverageRun announces that a coverage run is starting, stopping, or progressing.
	OnCoverageRun(rm RunMessage)
}

//go:generate mockery --name=Observer

// RunMessage is the type of messages announcing coverage run progress.
type RunMessage struct {
	observing.Batch

	// ProfileName contains the runner profile name on start and stop messages.
	ProfileName string

	// Profile contains the runner profile on start messages.
	Profile *Profile

	// Context contains the runner context on step messages.
	Context *RunContext
}

// OnCoverageRun broadcasts run message rm to all observers o.
func OnCoverageRun(rm RunMessage, o ...Observer) {
	for _, obs := range o {
		obs.OnCoverageRun(rm)
	}
}

// RunStart constructs a message stating that a coverage run of size nruns is starting using profile p (named pname).
func RunStart(pname string, p Profile, nruns int) RunMessage {
	return RunMessage{
		Batch: observing.Batch{
			Kind: observing.BatchStart,
			Num:  nruns,
		},
		ProfileName: pname,
		Profile:     &p,
	}
}

// RunStep announces that a coverage run instance, number i, in profile pname and with runner context rc, is starting.
func RunStep(pname string, i int, rc RunContext) RunMessage {
	return RunMessage{
		Batch: observing.Batch{
			Kind: observing.BatchStep,
			Num:  i,
		},
		ProfileName: pname,
		Context:     &rc,
	}
}

// RunEnd announces that a coverage run is ending for the profile named pname.
func RunEnd(pname string) RunMessage {
	return RunMessage{
		Batch: observing.Batch{
			Kind: observing.BatchEnd,
		},
		ProfileName: pname,
	}
}
