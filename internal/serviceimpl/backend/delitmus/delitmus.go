// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package delitmus lifts the ACT delitmusifier into a backend.
package delitmus

import (
	"context"
	"io"
	"path/filepath"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"

	"github.com/MattWindsor91/act-tester/internal/act"
	"github.com/MattWindsor91/act-tester/internal/model/job"
	"github.com/MattWindsor91/act-tester/internal/model/recipe"
	"github.com/MattWindsor91/act-tester/internal/model/service"
	"github.com/MattWindsor91/act-tester/internal/subject/obs"
)

// Delitmus partially implements the backend specification by delegating to ACT's delitmusifier.
//
// The delitmus backend can't actually produce standalone C code, and, at time of writing, there is no way to get
// the tester to compile C code without running it.  Instead, its main purpose is to serve as the target of a coverage
// run.
type Delitmus struct {
	// BaseRunner is the base configuration of the act runner, which is copied and overridden for each lifting.
	BaseRunner act.Runner
}

const (
	outAux = "aux.json"
	outC   = "delitmus.c"
)

// Lift delitmusifies the litmus file specified in j, using errw for standard output.
// It outputs a delitmusified C file and auxiliary file to j's output directory, and produces a recipe that suggests
// compiling that C file as an object.
// At time of writing, there is no way to specify how to delitmusify the file.
func (d Delitmus) Lift(ctx context.Context, j job.Lifter, errw io.Writer) (recipe.Recipe, error) {
	a := act.Runner{DuneExec: false, Stderr: errw}

	dj := act.DelitmusJob{
		InLitmus: j.In.Filepath(),
		OutAux:   filepath.Join(j.OutDir, outAux),
		OutC:     filepath.Join(j.OutDir, outC),
	}
	if err := a.Delitmus(ctx, dj); err != nil {
		return recipe.Recipe{}, err
	}
	return recipe.New(j.OutDir,
		recipe.AddFiles(dj.OutC),
		recipe.AddInstructions(recipe.CompileObjInst(1)),
		// TODO(@MattWindsor91): deal with the fact that this ends in an obj compile
	), nil
}

// ParseObs errors, for we cannot parse the observations of a delitmus run.
func (d Delitmus) ParseObs(_ context.Context, _ *service.Backend, _ io.Reader, _ *obs.Obs) error {
	return backend.ErrNotSupported
}
