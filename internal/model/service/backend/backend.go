// Copyright (c) 2020 Matt Windsor and contributors
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

	"github.com/c4-project/c4t/internal/model/id"
)

// Backend contains the various interfaces that a backend can implement.
type Backend interface {
	// Capabilities gets the capability set reported for this backend.
	Capabilities() Capability

	// LitmusArches gets the list of Litmus architectures understood by this backend (capability CanLiftLitmus).
	LitmusArches() []id.ID

	// Some backends can lift test-cases into recipes (capability CanLift).
	SingleLifter

	// Backends that can be run standalone or produce executables (capability CanRunStandalone | CanProduceExe)
	// must give an observation parser for interpreting their stdout as observations.
	ObsParser
}

// ErrNotSupported is the error that backends should return if we try to do something they don't support.
var ErrNotSupported = errors.New("service doesn't support action")

// BackendRunner is the interface that backends must implement to slot into the machine node runner.
type BackendRunner interface {
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

	perr := p.ParseObs(ctx, j.Backend, obsr, j.Obs)
	werr := cmd.Wait()
	return errhelp.FirstError(perr, werr)
}
*/

// Resolver is the interface of things that can resolve backends.
type Resolver interface {
	// Resolve tries to resolve the spec s into a backend.
	Resolve(s *Spec) (Backend, error)
}

//go:generate mockery --name=Resolver
