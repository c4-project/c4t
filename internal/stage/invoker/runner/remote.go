// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/MattWindsor91/c4t/internal/quantity"

	"github.com/MattWindsor91/c4t/internal/ux/stdflag"

	"github.com/MattWindsor91/c4t/internal/copier"

	"github.com/MattWindsor91/c4t/internal/plan"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/c4t/internal/remote"

	"github.com/alessio/shellescape"
	"golang.org/x/crypto/ssh"
)

// RemoteFactory is a factory that produces runners in the form of SSH sessions.
type RemoteFactory struct {
	// machine contains the instantiated machine runner, if present.
	machine *remote.MachineRunner
}

// NewRemoteFactory opens a SSH connection using Config and mc.
// If successful, it creates a runner factory over it.
func NewRemoteFactory(gc *remote.Config, mc *remote.MachineConfig) (*RemoteFactory, error) {
	machine, err := mc.MachineRunner(gc)
	return &RemoteFactory{machine: machine}, err
}

// MakeRunner constructs a runner using this factory's open SSH connection.
func (s *RemoteFactory) MakeRunner(ldir string, _ *plan.Plan, obs ...copier.Observer) (Runner, error) {
	// TODO(@MattWindsor91): re-establish connection if errors
	return NewRemoteRunner(s.machine, ldir, obs...)
}

// Close closes the underlying SSH connection being used for runners created by this factory.
func (s *RemoteFactory) Close() error {
	if s.machine == nil {
		return nil
	}
	return s.machine.Close()
}

// RemoteRunner runs the machine-runner via SSH.
type RemoteRunner struct {
	// observers observe any copying this RemoteRunner does.
	observers []copier.Observer
	// runner is the top-level runner to use for opening sessions and SFTP.
	runner *remote.MachineRunner
	// session receives the session once we start running the command.
	session *ssh.Session
	// localRoot is the slash-path of the root directory into which compile files should be received.
	localRoot string
	// remoteRoot is the slash-path of the remote directory into which compile files should be sent.
	remoteRoot string
	// eg is used to coordinate the combination of waiting for the SSH transaction to close and listening for the
	// context cancelling underneath it.
	eg errgroup.Group
}

// NewRemoteRunner creates a new RemoteRunner.
func NewRemoteRunner(r *remote.MachineRunner, localRoot string, o ...copier.Observer) (*RemoteRunner, error) {
	return &RemoteRunner{runner: r, observers: o, localRoot: localRoot, remoteRoot: r.Config.DirCopy}, nil
}

// Start starts a SSH session connected to a machine node with the quantities specified in qs.
func (r *RemoteRunner) Start(ctx context.Context, qs quantity.MachNodeSet) (*remote.Pipeset, error) {
	var (
		err error
		ps  *remote.Pipeset
	)

	if r.session, err = r.runner.NewSession(); err != nil {
		return nil, err
	}

	if ps, err = r.openPipes(); err != nil {
		return nil, fmt.Errorf("while opening pipes: %w", err)
	}

	if err := r.session.Start(r.invocation(qs)); err != nil {
		_ = ps.Close()
		return nil, fmt.Errorf("while starting local runner: %w", err)
	}

	makeSSHWaiter(&r.eg, r, ctx)

	return ps, nil
}

func makeSSHWaiter(eg *errgroup.Group, r *RemoteRunner, ctx context.Context) {
	// This channel makes sure that the context monitoring goroutine doesn't block.
	cl := make(chan struct{})
	*eg = errgroup.Group{}
	eg.Go(func() error {
		err := r.session.Wait()
		close(cl)
		return err
	})
	eg.Go(func() error {
		select {
		case <-cl:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
}

// Wait waits for either the SSH session to finish, or the context supplied to Start to close.
func (r *RemoteRunner) Wait() error {
	err := r.eg.Wait()

	// I'm unsure as to whether a session close errors if the session has been waited on;
	// hence why this error is currently unhandled.
	_ = r.session.Close()
	r.session = nil

	return err
}

// invocation works out what the SSH command invocation for the tester should be.
func (r *RemoteRunner) invocation(qs quantity.MachNodeSet) string {
	dir := path.Join(r.remoteRoot, "mach")
	qdir := shellescape.Quote(dir)
	return strings.Join(stdflag.MachInvocation(qdir, qs), " ")
}

// openPipes tries to open stdin, stdout, and stderr pipes for r.
func (r *RemoteRunner) openPipes() (*remote.Pipeset, error) {
	return remote.OpenSSHPipes(r.session)
}
