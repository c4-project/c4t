// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package herdtools

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/MattWindsor91/c4t/internal/helper/errhelp"

	backend2 "github.com/MattWindsor91/c4t/internal/model/service/backend"

	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend"

	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend/herdtools/parser"

	"github.com/MattWindsor91/c4t/internal/model/recipe"

	"github.com/MattWindsor91/c4t/internal/model/service"
	"github.com/MattWindsor91/c4t/internal/subject/obs"
)

// standaloneOut is the name of the file in the output directory to which we should write standalone output.
const standaloneOut = "output.txt"

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

func (h Backend) Lift(ctx context.Context, j backend2.LiftJob, x service.Runner) (recipe.Recipe, error) {
	if err := h.checkAndAmendJob(&j); err != nil {
		return recipe.Recipe{}, err
	}

	b := j.Backend
	if b == nil {
		return recipe.Recipe{}, fmt.Errorf("%w: backend in harness job", service.ErrNil)
	}

	r := h.DefaultRun
	if b.Run != nil {
		r.Override(*b.Run)
	}

	if err := h.liftInner(ctx, j, r, x); err != nil {
		return recipe.Recipe{}, fmt.Errorf("running %s: %w", r.Cmd, err)
	}
	return h.makeRecipe(j)
}

func (h Backend) liftInner(ctx context.Context, j backend2.LiftJob, r service.RunInfo, x service.Runner) error {
	var err error
	switch j.Out.Target {
	case backend2.ToStandalone:
		err = h.liftStandalone(ctx, j, r, x)
	case backend2.ToExeRecipe:
		err = h.Impl.LiftExe(ctx, j, r, x)
	}
	// We should've filtered out bad targets by this stage.
	return err
}

func (h Backend) liftStandalone(ctx context.Context, j backend2.LiftJob, r service.RunInfo, x service.Runner) error {
	f, err := os.Create(filepath.Join(filepath.Clean(j.Out.Dir), standaloneOut))
	if err != nil {
		return fmt.Errorf("couldn't create standalone output file: %s", err)
	}
	rerr := h.Impl.LiftStandalone(ctx, j, r, x, f)
	cerr := f.Close()
	return errhelp.FirstError(rerr, cerr)
}

func (h Backend) checkAndAmendJob(j *backend2.LiftJob) error {
	if err := j.Check(); err != nil {
		return err
	}
	if j.In.Source != backend2.LiftLitmus {
		return fmt.Errorf("%w: can only lift litmus files", backend.ErrNotSupported)
	}
	if j.Out.Target == backend2.ToDefault {
		j.Out.Target = backend2.ToStandalone
	}
	switch j.Out.Target {
	case backend2.ToStandalone:
	case backend2.ToExeRecipe:
		if (h.Capability & backend.CanProduceExe) == 0 {
			return fmt.Errorf("%w: cannot produce executables", backend.ErrNotSupported)
		}
	case backend2.ToObjRecipe:
		return fmt.Errorf("%w: cannot produce objects", backend.ErrNotSupported)
	}
	return nil
}

func (h Backend) makeRecipe(j backend2.LiftJob) (recipe.Recipe, error) {
	fs, err := j.Out.Files()
	if err != nil {
		return recipe.Recipe{}, err
	}

	return recipe.New(j.Out.Dir,
		targetRecipeOutput(j.Out.Target),
		recipe.AddFiles(fs...),
		// TODO(@MattWindsor91): delitmus support
		targetRecipeOption(j.Out.Target),
	)
}

func targetRecipeOutput(tgt backend2.Target) recipe.Output {
	if tgt == backend2.ToExeRecipe {
		return recipe.OutExe
	}
	return recipe.OutNothing
}

func targetRecipeOption(tgt backend2.Target) recipe.Option {
	if tgt == backend2.ToExeRecipe {
		return recipe.CompileAllCToExe()
	}
	return recipe.Options()
}

// BackendImpl describes the functionality that differs between Herdtools-style backends.
type BackendImpl interface {
	// LiftStandalone runs the lifter job j using x and the run information in r, expecting it to output the
	// results into w.
	LiftStandalone(ctx context.Context, j backend2.LiftJob, r service.RunInfo, x service.Runner, w io.Writer) error

	// LiftExe runs the lifter job j using x and the run information in r, expecting an executable.
	LiftExe(ctx context.Context, j backend2.LiftJob, r service.RunInfo, x service.Runner) error

	parser.Impl
}
