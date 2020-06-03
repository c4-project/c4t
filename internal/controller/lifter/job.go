// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package lifter

import (
	"context"
	"fmt"
	"io"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/model/service"

	"github.com/MattWindsor91/act-tester/internal/model/job"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Job is the type of per-subject lifter jobs.
type Job struct {
	// Arches is the list of architectures for which this job is responsible.
	Arches []id.ID

	// Backend is the backend that this job will use.
	Backend *service.Backend

	// Driver is the single-lift driver for this job.
	Driver SingleLifter

	// Paths is the path resolver for this job.
	Paths Pather

	// Stderr is the writer to which lifter stderr should be redirected.
	Stderr io.Writer

	// Normalise is the subject that we are trying to lift.
	Subject subject.Named

	// Rng is the random number generator to use for fuzz seeds.
	Rng *rand.Rand

	// ResCh is the channel onto which each fuzzed subject should be sent.
	ResCh chan<- builder.Request
}

// Lift performs this lifting job.
func (j *Job) Lift(ctx context.Context) error {
	if err := j.check(); err != nil {
		return err
	}

	// TODO(@MattWindsor91): this used to be a parallel loop, but was causing file exhaustion.
	// Ideally, we'd have a means of parallelism that doesn't inadvertently scale up like this.
	for _, a := range j.Arches {
		if err := j.liftArch(ctx, a); err != nil {
			return err
		}
	}
	return nil
}

// check does some basic checking on the Job before starting to run it.
func (j *Job) check() error {
	if j.Backend == nil {
		return ErrNoBackend
	}
	if j.Driver == nil {
		return ErrDriverNil
	}
	// It's ok for j.Stderr to be nil, as the SingleLifter is expected to deal with it.
	return nil
}

func (j *Job) liftArch(ctx context.Context, arch id.ID) error {
	dir, derr := j.Paths.Path(arch, j.Subject.Name)
	if derr != nil {
		return fmt.Errorf("when getting subject dir: %w", derr)
	}

	path, perr := j.Subject.BestLitmus()
	if perr != nil {
		return perr
	}

	spec := job.Lifter{
		Backend: j.Backend,
		Arch:    arch,
		InFile:  path,
		OutDir:  dir,
	}

	r, err := j.Driver.Lift(ctx, spec, j.Stderr)
	if err != nil {
		return fmt.Errorf("when lifting %s (arch %s): %w", j.Subject.Name, arch, err)
	}

	return builder.RecipeRequest(j.Subject.Name, arch, r).SendTo(ctx, j.ResCh)
}
