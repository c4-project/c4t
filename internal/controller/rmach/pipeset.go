// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package rmach

import "io"

// Pipeset groups together the three pipes of a Runner.
type Pipeset struct {
	// Stdin is the standard input pipe.
	Stdin io.WriteCloser
	// Stdout is the standard output pipe.
	Stdout io.ReadCloser
	// Stderr is the standard error
	Stderr io.ReadCloser
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
