// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package backend contains style-to-backend resolution.
package backend

import (
	"context"
	"errors"
	"io"

	"github.com/c4-project/c4t/internal/model/recipe"

	"github.com/c4-project/c4t/internal/model/service"
	"github.com/c4-project/c4t/internal/subject/obs"

	"github.com/c4-project/c4t/internal/id"
)

// Backend contains the various interfaces that a backend can implement.
type Backend interface {
	// SingleLifter captures that some backends can lift test-cases into recipes (capability CanLiftLitmus).
	SingleLifter

	// ObsParser captures that any backends that can be run standalone or produce executables
	// (capability CanRunStandalone | CanProduceExe) must give an parser for interpreting their stdout as observations.
	ObsParser

	// Class gets the class of this backend.
	Class() Class
}

// Class contains information about a style of backend.
type Class interface {
	// Metadata gets information about this type of backend.
	Metadata() Metadata

	// Instantiate instantiates a class, producing a backend.
	Instantiate(spec Spec) Backend

	// Probe probes the local system for specifications that can be used to produce backends of this class.
	// It takes a runner sr and context ctx, for any external programs, and a style ID style that the resolver will
	// resolve to one of these specifications.
	Probe(ctx context.Context, sr service.Runner, style id.ID) ([]NamedSpec, error)
}

// Metadata contains metadata for a backend archetype.
type Metadata struct {
	// Capabilities is the set of capability flags enabled on this backend.
	Capabilities Capability

	// LitmusArches is the list of Litmus architectures understood by this backend (capability CanLiftLitmus).
	LitmusArches []id.ID
}

// ErrNotSupported is the error that backends should return if we try to do something they don't support.
var ErrNotSupported = errors.New("service doesn't support action")

// Runner is the interface that backends must implement to slot into the machine node runner.
type Runner interface {
	// RunBackend runs the backend run job j.
	RunBackend(ctx context.Context, j *RunJob) error
}

// SingleLifter is an interface capturing the ability to lift single jobs into recipes.
type SingleLifter interface {
	// Lift performs the lifting described by j.
	// It returns a recipe describing the files (C files, header files, etc.) created and how to use them, or an error.
	// Any external service running should happen by sr.
	Lift(ctx context.Context, j LiftJob, sr service.Runner) (recipe.Recipe, error)
}

//go:generate mockery --name=SingleLifter

// ObsParser is the interface of things that can parse test outcomes.
type ObsParser interface {
	// ParseObs parses the observation in reader r into o according to the backend configuration in b.
	// The backend described by b must have been used to produce the testcase outputting r.
	ParseObs(ctx context.Context, r io.Reader, o *obs.Obs) error
}

/*
// RunExeAndParse runs the program described by r, parses its output with p, and emits the observations into j.
// It does not yet support the stubbing-out of the runner used.
func RunExeAndParse(ctx context.Context, j *RunJob, r service.RunInfo, p ObsParser) error {
	// TODO(@MattWindsor91): it'd be nice if this could be delegated to service.Runner, but quite complicated.

	cmd := exec.CommandContext(ctx, r.Cmd, r.Args...)
	obsr, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("while opening pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("while starting program: %w", err)
	}

	perr := p.ParseObs(ctx, obsr, j.Obs)
	werr := cmd.Wait()
	return errhelp.FirstError(perr, werr)
}
*/

// Resolver is the interface of things that can resolve backends.
type Resolver interface {
	// Resolve tries to resolve the spec s into a backend.
	Resolve(s Spec) (Backend, error)

	// Probe uses sr to probe for backend specifications on this machine.
	Probe(ctx context.Context, sr service.Runner) ([]NamedSpec, error)
}

//go:generate mockery --name=Resolver
