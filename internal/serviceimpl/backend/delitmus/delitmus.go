// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package delitmus lifts the c4f delitmusifier into a backend.
package delitmus

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/model/service"

	backend2 "github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/c4f"
	"github.com/c4-project/c4t/internal/model/recipe"
	"github.com/c4-project/c4t/internal/subject/obs"
)

const (
	outAux = "aux.json"
	outC   = "delitmus.c"
)

// Delitmus partially implements the backend specification by delegating to C4's delitmusifier.
//
// The delitmus backend can't actually produce standalone C code, and, at time of writing, there is no way to get
// the tester to compile C code without running it.  Instead, its main purpose is to serve as the target of a coverage
// run.
type Delitmus struct {
	// BaseRunner is the base configuration of the c4f runner, which is copied and overridden for each lifting.
	BaseRunner c4f.Runner
}

var dlMeta = backend2.Metadata{
	Capabilities: backend2.CanLiftLitmus | backend2.CanProduceObj,
	LitmusArches: []id.ID{id.ArchC},
}

// Metadata gets the metadata for the delitmusifier.
func (Delitmus) Metadata() backend2.Metadata {
	return dlMeta
}

// Instantiate 'instantiates' the delitmusifier; in fact, there isn't anything to instantiate.
func (d Delitmus) Instantiate(_ backend2.Spec) backend2.Backend {
	return d
}

// Class gets the 'class' of the delitmusifier (which is just the delitmusifier).
func (d Delitmus) Class() backend2.Class {
	return d
}

// Probe probes to see if there is a c4f installation we can use for delitmusifying.
func (Delitmus) Probe(ctx context.Context, sr service.Runner, style id.ID) ([]backend2.NamedSpec, error) {
	cr := c4f.Runner{DuneExec: false, Base: sr}

	// There's no actual information in the version flag yet.
	_, err := cr.CVersion(ctx)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	id.FromString("delitmus")
	return []backend2.NamedSpec{{ID: style, Spec: backend2.Spec{Style: style}}}, nil
}

// Lift delitmusifies the litmus file specified in j, using errw for standard output.
// It outputs a delitmusified C file and auxiliary file to j's output directory, and produces a recipe that suggests
// compiling that C file as an object.
// At time of writing, there is no way to specify how to delitmusify the file.
func (d Delitmus) Lift(ctx context.Context, j backend2.LiftJob, sr service.Runner) (recipe.Recipe, error) {
	if err := checkAndAmendJob(&j); err != nil {
		return recipe.Recipe{}, err
	}

	// Copying here is important; BaseRunner shouldn't have its service.Runner replaced
	a := d.BaseRunner
	a.Base = sr

	dj := c4f.DelitmusJob{
		InLitmus: j.In.Litmus.Filepath(),
		OutAux:   filepath.Join(j.Out.Dir, outAux),
		OutC:     filepath.Join(j.Out.Dir, outC),
	}
	if err := a.Delitmus(ctx, dj); err != nil {
		return recipe.Recipe{}, err
	}
	return recipe.New(j.Out.Dir,
		recipe.OutObj,
		recipe.AddFiles(dj.OutC),
		recipe.AddInstructions(recipe.CompileObjInst(1)),
	)
}

func checkAndAmendJob(j *backend2.LiftJob) error {
	if j.In.Source != backend2.LiftLitmus {
		return fmt.Errorf("%w: source must be litmus", backend2.ErrNotSupported)
	}
	if !j.In.Litmus.IsC() {
		return fmt.Errorf("%w: source must be C litmus", backend2.ErrNotSupported)
	}
	return checkAndAmendTarget(&j.Out.Target)
}

func checkAndAmendTarget(t *backend2.Target) error {
	switch *t {
	case backend2.ToDefault:
		*t = backend2.ToObjRecipe
		fallthrough
	case backend2.ToObjRecipe:
		return nil
	default:
		return fmt.Errorf("%w: cannot produce %q, only objects", backend2.ErrNotSupported, *t)
	}
}

// ParseObs errors, for we cannot parse the observations of a delitmus run.
func (d Delitmus) ParseObs(_ context.Context, _ io.Reader, _ *obs.Obs) error {
	return backend2.ErrNotSupported
}
