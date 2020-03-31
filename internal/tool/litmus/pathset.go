// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"errors"
	"path"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

var (
	// ErrNoFileIn occurs when the input file is empty.
	ErrNoFileIn = errors.New("need input file")

	// ErrNoDirOut occurs when the output directory is empty.
	ErrNoDirOut = errors.New("need output directory")
)

// Pathset contains the paths used in a litmus invocation.
type Pathset struct {
	// FileIn is the input file path.
	FileIn string

	// DirOut is the output directory path.
	DirOut string
}

// Check checks that the various paths in Pathset are populated.
// It doesn't check that they exist.
func (p *Pathset) Check() error {
	if p.FileIn == "" {
		return ErrNoFileIn
	}
	if p.DirOut == "" {
		return ErrNoDirOut
	}
	return nil
}

// Args gets this pathset's contribution to the litmus invocation arguments.
// These arguments should go at the end of the invocation.
func (p *Pathset) Args() []string {
	return []string{"-o", p.DirOut, p.FileIn}
}

// MainCFile guesses where Litmus is going to put the main file in its generated harness.
func (p *Pathset) MainCFile() string {
	file := iohelp.ExtlessFile(p.FileIn) + ".c"
	return path.Join(p.DirOut, file)
}
