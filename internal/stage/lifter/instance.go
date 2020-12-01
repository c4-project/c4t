// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package lifter

import (
	"context"
	"fmt"

	"github.com/MattWindsor91/c4t/internal/model/service"

	backend2 "github.com/MattWindsor91/c4t/internal/model/service/backend"

	"github.com/MattWindsor91/c4t/internal/model/id"

	"github.com/MattWindsor91/c4t/internal/subject/corpus/builder"

	"github.com/MattWindsor91/c4t/internal/subject"
)

// Instance is the type of per-subject lifter jobs.
type Instance struct {
	// Arches is the list of architectures for which this job is responsible.
	Arches []id.ID

	// Backend is the backend that this job will use.
	Backend *backend2.Spec

	// Driver is the single-lift driver for this job.
	Driver SingleLifter

	// Paths is the path resolver for this job.
	Paths Pather

	// Stderr is the runner on which the lifter should be run.
	Runner service.Runner

	// Normalise is the subject that we are trying to lift.
	Subject subject.Named

	// ResCh is the channel onto which each fuzzed subject should be sent.
	ResCh chan<- builder.Request
}

// Lift performs this lifting job.
func (j *Instance) Lift(ctx context.Context) error {
	if err := j.check(); err != nil {
		return err
	}

	// This used to be a parallel loop, but was contributing file exhaustion.  It might be safe to re-parallelise.
	for _, a := range j.Arches {
		if err := j.liftArch(ctx, a); err != nil {
			return err
		}
	}
	return nil
}

// check does some basic checking on the Instance before starting to run it.
func (j *Instance) check() error {
	if j.Backend == nil {
		return ErrNoBackend
	}
	if j.Driver == nil {
		return ErrDriverNil
	}
	// It's ok for j.Stderr to be nil, as the SingleLifter is expected to deal with it.
	return nil
}

func (j *Instance) liftArch(ctx context.Context, arch id.ID) error {
	dir, derr := j.Paths.Path(arch, j.Subject.Name)
	if derr != nil {
		return fmt.Errorf("when getting subject dir: %w", derr)
	}

	lit, perr := j.Subject.BestLitmus()
	if perr != nil {
		return perr
	}

	// TODO(@MattWindsor91): don't hardcode this
	spec := backend2.LiftJob{
		Backend: j.Backend,
		Arch:    arch,
		In:      backend2.LiftLitmusInput(*lit),
		Out: backend2.LiftOutput{
			Dir:    dir,
			Target: backend2.ToExeRecipe,
		},
	}

	r, err := j.Driver.Lift(ctx, spec, j.Runner)
	if err != nil {
		return fmt.Errorf("when lifting %s (arch %s): %w", j.Subject.Name, arch, err)
	}

	return builder.RecipeRequest(j.Subject.Name, arch, r).SendTo(ctx, j.ResCh)
}
