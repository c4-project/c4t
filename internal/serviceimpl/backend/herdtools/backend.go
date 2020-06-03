// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package herdtools

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/MattWindsor91/act-tester/internal/model/job"

	"github.com/MattWindsor91/act-tester/internal/model/obs"
	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// Backend represents Backend-style backends such as Herd and Litmus.
type Backend struct {
	// DefaultRun is the default run information for the particular backend.
	DefaultRun service.RunInfo

	// Impl provides parts of the Backend backend setup that differ between the various tools.
	Impl BackendImpl
}

// ParseObs parses an observation from r into o.
func (h Backend) ParseObs(_ context.Context, _ *service.Backend, r io.Reader, o *obs.Obs) error {
	p := parser{impl: h.Impl, o: o}
	s := bufio.NewScanner(r)
	lineno := 1
	for s.Scan() {
		if err := p.processLine(s.Text()); err != nil {
			return fmt.Errorf("line %d: %w", lineno, err)
		}
		lineno++
	}
	if err := s.Err(); err != nil {
		return err
	}
	return p.checkFinalState()
}

func (h Backend) Lift(ctx context.Context, j job.Lifter, errw io.Writer) (recipe.Recipe, error) {
	b := j.Backend

	if b == nil {
		return recipe.Recipe{}, fmt.Errorf("%w: backend in harness job", service.ErrNil)
	}

	r := h.DefaultRun
	if b.Run != nil {
		r.Override(*b.Run)
	}
	args, err := h.Impl.Args(j, r)
	if err != nil {
		return recipe.Recipe{}, err
	}

	cmd := exec.CommandContext(ctx, r.Cmd, args...)
	cmd.Stderr = errw
	if err := cmd.Run(); err != nil {
		return recipe.Recipe{}, fmt.Errorf("running %s: %w", r.Cmd, err)
	}

	return h.makeRecipe(j)
}

func (h Backend) makeRecipe(j job.Lifter) (recipe.Recipe, error) {
	fs, err := j.OutFiles()
	if err != nil {
		return recipe.Recipe{}, err
	}
	return recipe.New(j.OutDir,
		recipe.AddFiles(fs...),
		// TODO(@MattWindsor91): delitmus support
		recipe.CompileAllCToExe(),
	), nil
}

// BackendImpl describes the functionality that differs between Herdtools-style backends.
type BackendImpl interface {
	// Args tries to deduce the arguments needed to run the harness job j according to merged run information r.
	// It can fail if the job is not runnable by the backend.
	Args(j job.Lifter, r service.RunInfo) ([]string, error)

	// ParseStateCount parses the state-count line whose raw fields are fields.
	ParseStateCount(fields []string) (uint64, error)

	// ParseStateLine parses the state line whose raw fields are fields.
	ParseStateLine(tt TestType, fields []string) (*StateLine, error)
}
