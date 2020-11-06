// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package herdtools

import (
	"context"
	"fmt"
	"io"

	backend2 "github.com/MattWindsor91/act-tester/internal/model/service/backend"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/herdtools/parser"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/MattWindsor91/act-tester/internal/model/service"
	"github.com/MattWindsor91/act-tester/internal/subject/obs"
)

// Backend represents herdtools-style backends such as Herd and Litmus.
type Backend struct {
	// Capability contains the capability flags for this backend.
	Capability backend.Capability

	// DefaultRun is the default run information for the particular backend.
	DefaultRun service.RunInfo

	// Impl provides parts of the Backend backend setup that differ between the various tools.
	Impl BackendImpl
}

// Capabilities returns Capability, to satisfy the backend interface.
func (h Backend) Capabilities(_ *backend2.Spec) backend.Capability {
	return h.Capability
}

// ParseObs parses an observation from r into o.
func (h Backend) ParseObs(_ context.Context, _ *backend2.Spec, r io.Reader, o *obs.Obs) error {
	return parser.Parse(h.Impl, r, o)
}

func (h Backend) Lift(ctx context.Context, j backend2.LiftJob, sr service.Runner) (recipe.Recipe, error) {
	b := j.Backend

	if b == nil {
		return recipe.Recipe{}, fmt.Errorf("%w: backend in harness job", service.ErrNil)
	}

	r := h.DefaultRun
	if b.Run != nil {
		r.Override(*b.Run)
	}

	if err := h.Impl.Run(ctx, j, r, sr); err != nil {
		return recipe.Recipe{}, fmt.Errorf("running %s: %w", r.Cmd, err)
	}

	return h.makeRecipe(j)
}

func (h Backend) makeRecipe(j backend2.LiftJob) (recipe.Recipe, error) {
	fs, err := j.Out.Files()
	if err != nil {
		return recipe.Recipe{}, err
	}
	return recipe.New(j.Out.Dir,
		recipe.AddFiles(fs...),
		// TODO(@MattWindsor91): delitmus support
		recipe.CompileAllCToExe(),
	), nil
}

// BackendImpl describes the functionality that differs between Herdtools-style backends.
type BackendImpl interface {
	// Run runs the lifter job j using x and the run information in r.
	Run(ctx context.Context, j backend2.LiftJob, r service.RunInfo, x service.Runner) error

	parser.Impl
}
