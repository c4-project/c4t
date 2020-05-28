// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package rmach

import (
	"context"
	"fmt"
	"path"
	"strings"

	copy2 "github.com/MattWindsor91/act-tester/internal/copier"

	"github.com/MattWindsor91/act-tester/internal/model/plan"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/remote"

	"github.com/alessio/shellescape"
	"golang.org/x/crypto/ssh"
)

// PlanRunnerFactory is a runner factory that instantiates either a SSH or local runner depending on the machine
// configuration inside the first plan passed to its MakeRunner.
type PlanRunnerFactory struct {
	recvRoot string
	gc       *remote.Config

	cached RunnerFactory
}

// MakeRunner makes a runner using the machine configuration in pl.
func (p *PlanRunnerFactory) MakeRunner(pl *plan.Plan, obs ...copy2.Observer) (Runner, error) {
	var err error
	if p.cached == nil {
		if p.cached, err = p.makeFactory(pl); err != nil {
			return nil, err
		}
	}
	return p.cached.MakeRunner(pl, obs...)
}

func (p *PlanRunnerFactory) makeFactory(pl *plan.Plan) (RunnerFactory, error) {
	if pl.Machine.SSH == nil {
		return LocalRunnerFactory(p.recvRoot), nil
	}
	return NewSSHRunnerFactory(p.recvRoot, p.gc, pl.Machine.SSH)
}

// Close closes the runner factory, if it was ever instantiated.
func (p *PlanRunnerFactory) Close() error {
	if p.cached == nil {
		return nil
	}
	return p.cached.Close()
}

// SSHRunnerFactory is a factory that produces runners in the form of SSH sessions.
type SSHRunnerFactory struct {
	recvRoot string
	// machine contains the instantiated machine runner, if present.
	machine *remote.MachineRunner
}

// NewSSHRunnerFactory opens a SSH connection using gc and mc.
// If successful, it creates a runner factory over it, using recvRoot as the local directory.
func NewSSHRunnerFactory(recvRoot string, gc *remote.Config, mc *remote.MachineConfig) (*SSHRunnerFactory, error) {
	machine, err := mc.MachineRunner(gc)
	return &SSHRunnerFactory{recvRoot: recvRoot, machine: machine}, err
}

// MakeRunner constructs a runner using this factory's open SSH connection.
func (s *SSHRunnerFactory) MakeRunner(_ *plan.Plan, obs ...copy2.Observer) (Runner, error) {
	// TODO(@MattWindsor91): re-establish connection if errors
	return NewSSHRunner(s.machine, s.recvRoot, obs...)
}

// Close closes the underlying SSH connection being used for runners created by this factory.
func (s *SSHRunnerFactory) Close() error {
	return s.machine.Close()
}

// SSHRunner runs the machine-runner via SSH.
type SSHRunner struct {
	// observers observe any copying this SSHRunner does.
	observers []copy2.Observer
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

// NewSSHRunner creates a new SSHRunner.
func NewSSHRunner(r *remote.MachineRunner, localRoot string, o ...copy2.Observer) (*SSHRunner, error) {
	return &SSHRunner{runner: r, observers: o, localRoot: localRoot, remoteRoot: r.Config.DirCopy}, nil
}

// Start starts a SSH session connected to a machine node with name and arguments constructed through i.
func (r *SSHRunner) Start(ctx context.Context, i InvocationGetter) (*remote.Pipeset, error) {
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

	if err := r.session.Start(r.invocation(i)); err != nil {
		_ = ps.Close()
		return nil, fmt.Errorf("while starting local runner: %w", err)
	}

	makeSSHWaiter(&r.eg, r, ctx)

	return ps, nil
}

func makeSSHWaiter(eg *errgroup.Group, r *SSHRunner, ctx context.Context) {
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
func (r *SSHRunner) Wait() error {
	err := r.eg.Wait()

	// I'm unsure as to whether a session close errors if the session has been waited on;
	// hence why this error is currently unhandled.
	_ = r.session.Close()
	r.session = nil

	return err
}

// invocation works out what the SSH command invocation for the tester should be.
func (r *SSHRunner) invocation(i InvocationGetter) string {
	dir := path.Join(r.remoteRoot, "mach")
	qdir := shellescape.Quote(dir)
	return strings.Join(Invocation(i, qdir), " ")
}

// openPipes tries to open stdin, stdout, and stderr pipes for r.
func (r *SSHRunner) openPipes() (*remote.Pipeset, error) {
	return remote.OpenSSHPipes(r.session)
}
