// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/MattWindsor91/c4t/internal/model/service"
)

// Runner is the interface of things that can be run to generate coverage testbeds.
type Runner interface {
	// Run runs the Runner with context ctx and runner context rc.
	Run(ctx context.Context, rc RunContext) error
}

//go:generate mockery --name=Runner

// StandaloneRunner is a coverage runner that runs a standalone binary.
type StandaloneRunner struct {
	// run tells the runner how to run the standalone runner.
	run service.RunInfo
	// errw is the writer to which stderr should go, if any.
	errw io.Writer
}

// Run runs the standalone runner.
func (s *StandaloneRunner) Run(ctx context.Context, rc RunContext) error {
	cmd := exec.CommandContext(ctx, s.run.Cmd, rc.ExpandArgs(s.run.Args...)...)
	cmd.Stderr = s.errw
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running coverage generator %q: %w", s.run.Cmd, err)
	}
	return nil
}
