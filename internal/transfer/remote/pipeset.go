// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package remote

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

// Pipeset groups together the three pipes of a Runner.
type Pipeset struct {
	// Stdin is the standard input pipe.
	Stdin io.WriteCloser
	// Stdout is the standard output pipe.
	Stdout io.ReadCloser
	// Stderr is the standard error
	Stderr io.ReadCloser
}

// OpenCmdPipes opens a pipeset on command c.
func OpenCmdPipes(c *exec.Cmd) (*Pipeset, error) {
	var (
		ps  Pipeset
		err error
	)
	if ps.Stdin, err = c.StdinPipe(); err != nil {
		return nil, fmt.Errorf("while opening stdin pipe: %w", err)
	}
	if ps.Stdout, err = c.StdoutPipe(); err != nil {
		_ = ps.Close()
		return nil, fmt.Errorf("while opening stdout pipe: %w", err)
	}
	if ps.Stderr, err = c.StderrPipe(); err != nil {
		_ = ps.Close()
		return nil, fmt.Errorf("while opening stderr pipe: %w", err)
	}
	return &ps, nil
}

// OpenSSHPipes opens a pipeset on SSH session s.
func OpenSSHPipes(s *ssh.Session) (*Pipeset, error) {
	stdin, err := s.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("while opening stdin pipe: %w", err)
	}
	stdout, err := s.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		return nil, fmt.Errorf("while opening stdout pipe: %w", err)
	}
	stderr, err := s.StderrPipe()
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

// Close tries to close each non-nil pipe in Pipeset.
func (p *Pipeset) Close() error {
	if err := safeClose(p.Stdin); err != nil {
		return err
	}
	if err := safeClose(p.Stdout); err != nil {
		return err
	}
	return safeClose(p.Stderr)
}

// safeClose closes c if, and only if, it is non-nil.
func safeClose(c io.Closer) error {
	if c == nil {
		return nil
	}
	return c.Close()
}
