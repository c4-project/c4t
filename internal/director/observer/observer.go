// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package observer defines interfaces and basic implementations of the director's 'observer' pattern.
package observer

import (
	"context"
	"io"

	copy2 "github.com/MattWindsor91/act-tester/internal/copier"

	"github.com/MattWindsor91/act-tester/internal/model/machine"

	"github.com/MattWindsor91/act-tester/internal/stage/analyse/observer"

	"github.com/MattWindsor91/act-tester/internal/model/run"

	"github.com/MattWindsor91/act-tester/internal/stage/planner"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// Observer is an interface for types that implement multi-machine test progress observation.
//
// Unlike most observer patterns in the tester, the director takes a degree of control over the lifecycle of its
// observers.  It will call Run for each observer in parallel with the tester instances, and call Close when its
// observers are no longer needed.
type Observer interface {
	machine.Observer

	// Run runs any runtime required by the observer in a blocking manner using context ctx.
	// It will use cancel to cancel ctx if needed.
	Run(ctx context.Context, cancel func()) error

	// The director will Close observers when it is shutting down.
	io.Closer

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

	// Instance observers can observe plan analyses.
	observer.Observer

	// Instance observers can observe planner operations.
	planner.Observer

	// Instance observers can observe file copies.
	copy2.Observer
}

// OnIteration sends OnIteration to every instance observer in obs.
func OnIteration(r run.Run, obs ...Instance) {
	for _, o := range obs {
		o.OnIteration(r)
	}
}

// LowerToMachine lowers a slice of director observers to a slice of machine observers.
func LowerToMachine(obs []Observer) []machine.Observer {
	cos := make([]machine.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToAnalyse lowers a slice of instance observers to a slice of builder observers.
func LowerToAnalyse(obs []Instance) []observer.Observer {
	cos := make([]observer.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
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
func LowerToCopy(obs []Instance) []copy2.Observer {
	cos := make([]copy2.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// CloseAll closes all observers passed to it, returning the error of the last one (if any).
func CloseAll(obs ...Observer) error {
	var err error
	for _, o := range obs {
		err = o.Close()
	}
	return err
}
