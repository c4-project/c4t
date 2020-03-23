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

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/pkg/director/observer"

	"github.com/MattWindsor91/act-tester/internal/pkg/transfer"

	"github.com/pkg/sftp"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/transfer/remote"

	"github.com/alessio/shellescape"
	"golang.org/x/crypto/ssh"
)

// SSHRunner runs the machine-runner via SSH.
type SSHRunner struct {
	// observer observes any copying this SSHRunner does.
	observer observer.Copy
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
func NewSSHRunner(r *remote.MachineRunner, o observer.Copy) *SSHRunner {
	// TODO(@MattWindsor91): recvRoot
	return &SSHRunner{runner: r, observer: o}
}

func (r *SSHRunner) Start(ctx context.Context) (*Pipeset, error) {
	// TODO(@MattWindsor91): handle context
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

// Send normalises p against the remote directory, then SFTPs each affected file into place on the remote machine.
func (r *SSHRunner) Send(p *plan.Plan) (*plan.Plan, error) {
	n := transfer.NewNormaliser(r.runner.Config.DirCopy)
	rp := *p
	var err error
	if rp.Corpus, err = n.Corpus(rp.Corpus); err != nil {
		return nil, err
	}

	return &rp, r.sftpMappings(n.HarnessMappings())
}

func (r *SSHRunner) Recv(locp, remp *plan.Plan) (*plan.Plan, error) {
	for n, rs := range remp.Corpus {
		norm := transfer.NewNormaliser(path.Join(r.recvRoot, n))
		ls, ok := locp.Corpus[n]
		if !ok {
			return nil, fmt.Errorf("subject not in local corpus: %s", n)
		}
		ns, err := norm.Subject(rs)
		if err != nil {
			return nil, err
		}
		ls.Runs = ns.Runs
		// TODO(@MattWindsor91): receive
		locp.Corpus[n] = ls
	}

	return locp, nil
}

func (r *SSHRunner) sftpMappings(ms map[string]string) error {
	r.observer.OnCopyStart(len(ms))
	defer r.observer.OnCopyFinish()

	cli, err := r.runner.NewSFTP()
	if err != nil {
		return err
	}
	for rpath, lpath := range ms {
		if err := sftpMapping(cli, rpath, lpath); err != nil {
			_ = cli.Close()
			return err
		}
		r.observer.OnCopy(lpath, rpath)
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
