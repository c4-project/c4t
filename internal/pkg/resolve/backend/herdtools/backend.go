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

	"github.com/MattWindsor91/act-tester/internal/pkg/model/job"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/obs"
	"github.com/MattWindsor91/act-tester/internal/pkg/model/service"
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

func (h Backend) MakeHarness(ctx context.Context, j job.Harness, errw io.Writer) (outFiles []string, err error) {
	b := j.Backend
	if b == nil {
		return nil, fmt.Errorf("%w: backend in harness job", service.ErrNil)
	}

	r := h.DefaultRun
	if b.Run != nil {
		r.Override(*b.Run)
	}
	args, err := h.Impl.Args(j, r)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, r.Cmd, args...)
	cmd.Stderr = errw
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("running %s: %w", r.Cmd, err)
	}

	return j.OutFiles()
}

// BackendImpl describes the functionality that differs between Herdtools-style backends.
type BackendImpl interface {
	// Args tries to deduce the arguments needed to run the harness job j according to merged run information r.
	// It can fail if the job is not runnable by the backend.
	Args(j job.Harness, r service.RunInfo) ([]string, error)

	// ParseStateCount parses the state-count line whose raw fields are fields.
	ParseStateCount(fields []string) (uint64, error)

	// ParseStateLine parses the state line whose raw fields are fields.
	ParseStateLine(tt TestType, fields []string) (*StateLine, error)
}
