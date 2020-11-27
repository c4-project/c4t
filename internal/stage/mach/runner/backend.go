// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/MattWindsor91/c4t/internal/helper/errhelp"

	"github.com/MattWindsor91/c4t/internal/model/service"

	"github.com/MattWindsor91/c4t/internal/model/service/backend"
	backend2 "github.com/MattWindsor91/c4t/internal/model/service/backend"
	"github.com/MattWindsor91/c4t/internal/subject/obs"
)

// BackendRunner is the interface that backends must implement to slot into the machine node runner.
type BackendRunner interface {
	// RunBackend runs the backend run job j.
	RunBackend(ctx context.Context, j *backend.RunJob) error
}

// ObsParser is the interface of things that can parse test outcomes.
type ObsParser interface {
	// ParseObs parses the observation in reader r into o according to the backend configuration in b.
	// The backend described by b must have been used to produce the testcase outputting r.
	ParseObs(ctx context.Context, b *backend2.Spec, r io.Reader, o *obs.Obs) error
}

// RunExeAndParse runs the program described by r, parses its output with p, and emits the observations into j.
// It does not yet support the stubbing-out of the runner used.
func RunExeAndParse(ctx context.Context, j *backend.RunJob, r service.RunInfo, p ObsParser) error {
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
