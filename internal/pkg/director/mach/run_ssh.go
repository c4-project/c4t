// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package mach

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/pkg/remote"

	"github.com/alessio/shellescape"
	"golang.org/x/crypto/ssh"
)

// SSHRunner runs the machine-runner via SSH.
type SSHRunner struct {
	// runner tells us how to run SSH.
	runner *remote.MachineRunner
	// session receives the session once we start running the command.
	session *ssh.Session
}

// NewSSHRunner creates a new SSHRunner.
func NewSSHRunner(r *remote.MachineRunner) *SSHRunner {
	return &SSHRunner{runner: r}
}

func (r *SSHRunner) Start(_ context.Context) (*Pipeset, error) {
	// TODO(@MattWindsor): handle context
	var err error

	r.session, err = r.runner.NewSession()
	if err != nil {
		return nil, err
	}

	var ps *Pipeset
	if ps, err = r.openPipes(); err != nil {
		return nil, fmt.Errorf("while opening pipes: %w", err)
	}

	if err := r.session.Start(r.invocation()); err != nil {
		_ = ps.Close()
		return nil, fmt.Errorf("while starting local runner: %w", err)
	}

	return ps, nil
}

func (r *SSHRunner) Wait() error {
	return r.session.Wait()
}

// invocation works out what the SSH command invocation for the tester should be.
func (r *SSHRunner) invocation() string {
	dir := path.Join(r.runner.Config.DirCopy, "mach")
	qdir := shellescape.Quote(dir)
	return strings.Join(runArgs(qdir), " ")
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
