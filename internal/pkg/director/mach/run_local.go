// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package mach

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/plan"
)

// LocalRunner runs the machine-runner locally.
type LocalRunner struct {
	// dir is the directory in which we are running the machine-runner.
	dir string
	// cmd receives the command once we start running the LocalRunner.
	cmd *exec.Cmd
}

// NewLocalRunner creates a new LocalRunner.
func NewLocalRunner(dir string) *LocalRunner {
	return &LocalRunner{dir: dir}
}

// Start starts the machine-runner binary locally using ctx, and returns various
func (r *LocalRunner) Start(ctx context.Context) (*Pipeset, error) {
	r.cmd = exec.CommandContext(ctx, binName, runArgs(r.dir)...)
	ps, err := r.openPipes()
	if err != nil {
		return nil, fmt.Errorf("opening pipes: %w", err)
	}
	err = r.cmd.Start()
	if err != nil {
		_ = ps.Close()
		return nil, fmt.Errorf("starting command: %w", err)
	}
	return ps, nil
}

// Send effectively does nothing but implement the general runner interface obligations.
func (r *LocalRunner) Send(p *plan.Plan) (*plan.Plan, error) {
	return p, nil
}

// Wait waits for the running machine-runner binary to terminate.
func (r *LocalRunner) Wait() error {
	return r.cmd.Wait()
}

// Recv effectively does nothing but implement the general runner interface obligations.
func (r *LocalRunner) Recv(_, rp *plan.Plan) (*plan.Plan, error) {
	// rp has been created on the local machine without any modifications, and needs no merging into the local plan.
	return rp, nil
}

// openLocalPipes tries to open stdin, stdout, and stderr pipes for c.
func (r *LocalRunner) openPipes() (*Pipeset, error) {
	var (
		ps  Pipeset
		err error
	)
	if ps.Stdin, err = r.cmd.StdinPipe(); err != nil {
		return nil, fmt.Errorf("while opening stdin pipe: %w", err)
	}
	if ps.Stdout, err = r.cmd.StdoutPipe(); err != nil {
		_ = ps.Close()
		return nil, fmt.Errorf("while opening stdout pipe: %w", err)
	}
	if ps.Stderr, err = r.cmd.StderrPipe(); err != nil {
		_ = ps.Close()
		return nil, fmt.Errorf("while opening stderr pipe: %w", err)
	}
	return &ps, nil
}
