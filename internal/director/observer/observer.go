// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package observer defines interfaces and basic implementations of the director's 'observer' pattern.
package observer

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/model/run"

	"github.com/MattWindsor91/act-tester/internal/controller/planner"

	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"

	"github.com/MattWindsor91/act-tester/internal/remote"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// Observer is an interface for types that implement multi-machine test progress observation.
type Observer interface {
	// Run runs the observer in a blocking manner using context ctx.
	// It will use cancel to cancel ctx if needed.
	Run(ctx context.Context, cancel func()) error

	// Instance gets a sub-observer for the machine with ID id.
	// It can fail if no such observer is available.
	Instance(id id.ID) (Instance, error)
}

// Instance is an interface for types that observe a director loop.
type Instance interface {
	// OnIteration lets the observer know that the machine loop has started anew.
	// iter is, modulo eventual overflow, the current iteration number;
	// time is the time at which the iteration started.
	OnIteration(run run.Run)

	// OnCollation lets the observer know that the run results have been received and collated into c.
	OnCollation(c *analysis.Analysis)

	// Instance observers can observe planner operations.
	planner.Observer

	// Instance observers can observe file copies.
	remote.CopyObserver
}

// OnIteration sends OnIteration to every instance observer in obs.
func OnIteration(r run.Run, obs ...Instance) {
	for _, o := range obs {
		o.OnIteration(r)
	}
}

// OnCollation sends OnCollation to every instance observer in obs.
func OnCollation(c *analysis.Analysis, obs ...Instance) {
	for _, o := range obs {
		o.OnCollation(c)
	}
}

// LowerToPlanner lowers a slice of instance observers to a slice of planner observers.
func LowerToPlanner(obs []Instance) []planner.Observer {
	cos := make([]planner.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToBuilder lowers a slice of instance observers to a slice of builder observers.
func LowerToBuilder(obs []Instance) []builder.Observer {
	cos := make([]builder.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToCopy lowers a slice of instance observers to a slice of copy observers.
func LowerToCopy(obs []Instance) []remote.CopyObserver {
	cos := make([]remote.CopyObserver, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}
