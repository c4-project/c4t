// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package observer defines interfaces and basic implementations of the director's 'observer' pattern.
package observer

import (
	"context"
	"io"

	"github.com/MattWindsor91/act-tester/internal/copier"

	mach "github.com/MattWindsor91/act-tester/internal/stage/mach/observer"

	"github.com/MattWindsor91/act-tester/internal/stage/perturber"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser"

	"github.com/MattWindsor91/act-tester/internal/machine"

	"github.com/MattWindsor91/act-tester/internal/model/run"

	"github.com/MattWindsor91/act-tester/internal/stage/planner"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"
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
	analyser.Observer

	// Instance observers can observe plan saves.
	saver.Observer

	// Instance observers can observe perturber operations.
	perturber.Observer

	// Instance observers can observe planner operations.
	planner.Observer

	// Instance observers can observe file copies.
	copier.Observer

	// Instance observers can observe machine node actions.
	mach.Observer
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

// LowerToAnalyser lowers a slice of instance observers to a slice of analyser observers.
func LowerToAnalyser(obs []Instance) []analyser.Observer {
	cos := make([]analyser.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToPerturber lowers a slice of instance observers to a slice of perturber observers.
func LowerToPerturber(obs []Instance) []perturber.Observer {
	cos := make([]perturber.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToSave lowers a slice of instance observers to a slice of saver observers.
func LowerToSaver(obs []Instance) []saver.Observer {
	cos := make([]saver.Observer, len(obs))
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
func LowerToCopy(obs []Instance) []copier.Observer {
	cos := make([]copier.Observer, len(obs))
	for i, o := range obs {
		cos[i] = o
	}
	return cos
}

// LowerToMach lowers a slice of director observers to a slice of machine node observers.
func LowerToMach(obs []Instance) []mach.Observer {
	cos := make([]mach.Observer, len(obs))
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
