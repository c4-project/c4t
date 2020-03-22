// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package mach

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/pkg/transfer"

	"github.com/pkg/sftp"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/transfer/remote"

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

// Send normalises p against the remote directory, then SFTPs each affected file into place on the remote machine.
func (r *SSHRunner) Send(p *plan.Plan) (*plan.Plan, error) {
	n := transfer.NewNormaliser(r.runner.Config.DirCopy)
	rp := *p
	var err error
	if rp.Corpus, err = n.Corpus(rp.Corpus); err != nil {
		return nil, err
	}

	return &rp, r.sftpMappings(n.Mappings)
}

func (r *SSHRunner) sftpMappings(ms map[string]string) error {
	cli, err := r.runner.NewSFTP()
	if err != nil {
		return err
	}
	for rpath, lpath := range ms {
		if err := sftpMapping(cli, rpath, lpath); err != nil {
			_ = cli.Close()
			return err
		}
	}
	return cli.Close()
}

func sftpMapping(cli *sftp.Client, rpath string, lpath string) error {
	if err := cli.MkdirAll(path.Dir(rpath)); err != nil {
		return err
	}

	r, err := os.Open(filepath.FromSlash(lpath))
	if err != nil {
		return err
	}
	w, err := cli.Create(rpath)
	if err != nil {
		_ = r.Close()
		return err
	}

	_, cperr := io.Copy(w, r)
	wcerr := w.Close()
	rcerr := r.Close()

	if cperr != nil {
		return cperr
	}
	if wcerr != nil {
		return wcerr
	}
	return rcerr
}

func (r *SSHRunner) Wait() error {
	return r.session.Wait()
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
