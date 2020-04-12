// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package rmach

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/transfer/remote"

	"github.com/alessio/shellescape"
	"golang.org/x/crypto/ssh"
)

// SSHRunner runs the machine-runner via SSH.
type SSHRunner struct {
	// observers observe any copying this SSHRunner does.
	observers []remote.CopyObserver
	// runner tells us how to run SSH.
	runner *remote.MachineRunner
	// session receives the session once we start running the command.
	session *ssh.Session
	// recvRoot is the slash-path of the root directory into which compile files should be received.
	recvRoot string
	// eg is used to coordinate the combination of waiting for the SSH transaction to close and listening for the
	// context cancelling underneath it.
	eg errgroup.Group
}

// NewSSHRunner creates a new SSHRunner.
func NewSSHRunner(r *remote.MachineRunner, o []remote.CopyObserver, recvRoot string) *SSHRunner {
	return &SSHRunner{runner: r, observers: o, recvRoot: recvRoot}
}

func (r *SSHRunner) Start(ctx context.Context) (*Pipeset, error) {
	var (
		err error
		ps  *Pipeset
	)

	if r.session, err = r.runner.NewSession(); err != nil {
		return nil, err
	}

	if ps, err = r.openPipes(); err != nil {
		return nil, fmt.Errorf("while opening pipes: %w", err)
	}

	if err := r.session.Start(r.invocation()); err != nil {
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
	return r.eg.Wait()
}

// invocation works out what the SSH command invocation for the tester should be.
func (r *SSHRunner) invocation() string {
	dir := path.Join(r.runner.Config.DirCopy, "mach")
	qdir := shellescape.Quote(dir)
	argv := append([]string{binName}, runArgs(qdir)...)
	return strings.Join(argv, " ")
}

// openPipes tries to open stdin, stdout, and stderr pipes for c.
func (r *SSHRunner) openPipes() (*Pipeset, error) {
	stdin, err := r.session.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("while opening stdin pipe: %w", err)
	}
	stdout, err := r.session.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		return nil, fmt.Errorf("while opening stdout pipe: %w", err)
	}
	stderr, err := r.session.StderrPipe()
	if err != nil {
		_ = stdin.Close()
		return nil, fmt.Errorf("while opening stderr pipe: %w", err)
	}
	ps := Pipeset{
		Stdin:  stdin,
		Stdout: ioutil.NopCloser(stdout),
		Stderr: ioutil.NopCloser(stderr),
	}
	return &ps, nil
}
