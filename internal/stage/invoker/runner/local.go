// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/ux/stdflag"

	copy2 "github.com/c4-project/c4t/internal/copier"

	"github.com/c4-project/c4t/internal/remote"

	"github.com/c4-project/c4t/internal/helper/iohelp"

	"github.com/c4-project/c4t/internal/plan"
)

// LocalFactory allows spawning of local runners using said path as the local directory.
type LocalFactory struct{}

func (l LocalFactory) MakeRunner(ldir string, _ *plan.Plan, _ ...copy2.Observer) (Runner, error) {
	return NewLocalRunner(ldir), nil
}

// Close does nothing.
func (l LocalFactory) Close() error { return nil }

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

// Start starts the machine-runner binary locally using ctx, and returns a pipeset for talking to it.
func (r *LocalRunner) Start(ctx context.Context, qs quantity.MachNodeSet) (*remote.Pipeset, error) {
	r.cmd = exec.CommandContext(ctx, stdflag.MachBinName, stdflag.MachArgs(r.dir, qs)...)
	ps, err := r.openPipes()
	if err != nil {
		return nil, fmt.Errorf("opening pipes: %w", err)
	}
	if err = r.cmd.Start(); err != nil {
		_ = ps.Close()
		return nil, fmt.Errorf("starting command: %w", err)
	}
	return ps, nil
}

// Send effectively does nothing but implement the general runner interface obligations and make sure the context is live.
func (r *LocalRunner) Send(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	return p, iohelp.CheckDone(ctx)
}

// Wait waits for the running machine-runner binary to terminate.
func (r *LocalRunner) Wait() error {
	return r.cmd.Wait()
}

// Recv effectively does nothing but implement the general runner interface obligations and make sure the context is live.
func (r *LocalRunner) Recv(ctx context.Context, _, rp *plan.Plan) (*plan.Plan, error) {
	// rp has been created on the local machine without any modifications, and needs no merging into the local plan.
	return rp, iohelp.CheckDone(ctx)
}

// openLocalPipes tries to open stdin, stdout, and stderr pipes for c.
func (r *LocalRunner) openPipes() (*remote.Pipeset, error) {
	return remote.OpenCmdPipes(r.cmd)
}
