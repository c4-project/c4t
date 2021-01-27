// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package lifter

import (
	"context"
	"fmt"
	"reflect"

	"github.com/c4-project/c4t/internal/model/service"

	"github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/model/id"

	"github.com/c4-project/c4t/internal/subject/corpus/builder"

	"github.com/c4-project/c4t/internal/subject"
)

// Instance is the type of per-subject lifter jobs.
type Instance struct {
	// Arches is the list of architectures for which this job is responsible.
	Arches []id.ID

	// Driver is the single-lift driver for this job.
	Driver backend.SingleLifter

	// Paths is the path resolver for this job.
	Paths Pather

	// Runner is the runner on which the lifter should be run.
	Runner service.Runner

	// Subject is the subject that we are trying to lift.
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
	spec := backend.LiftJob{
		Arch: arch,
		In:   backend.LiftLitmusInput(lit),
		Out: backend.LiftOutput{
			Dir:    dir,
			Target: backend.ToExeRecipe,
		},
	}

	r, err := j.Driver.Lift(ctx, spec, j.Runner)
	if err != nil {
		sname := reflect.TypeOf(j.Driver).Name()
		return &Error{Subject: &j.Subject, ServiceName: sname, Job: spec, Inner: err}
	}

	return builder.RecipeRequest(j.Subject.Name, arch, r).SendTo(ctx, j.ResCh)
}

// Error contains an error that occurred while lifting, as well as context.
type Error struct {
	// ServiceName is a guess at the name of the service.
	ServiceName string

	// Subject is, if non-nil, the name of the subject being lifted.
	Subject *subject.Named

	// Job is the job that was being processed when the error occurred.
	Job backend.LiftJob
	// Spec is
	// Inner is the inner error.
	Inner error
}

// SubjectName gets any subject name for this error.
func (e *Error) SubjectName() string {
	if e.Subject == nil {
		return "(unknown)"
	}
	return e.Subject.Name
}

// Error gets the error string for this error.
func (e *Error) Error() string {
	return fmt.Sprintf("when lifting %s with %s (arch %s): %s", e.SubjectName(), e.ServiceName, e.Job.Arch, e.Inner)
}

func (e *Error) Unwrap() error {
	return e.Inner
}
