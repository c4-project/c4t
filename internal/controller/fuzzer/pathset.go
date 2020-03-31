// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"path"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

const (
	// segLitmus is the directory element added to the root directory to form the litmus directory.
	segLitmus = "litmus"

	// segTrace is the directory element added to the root directory to form the trace directory.
	segTrace = "trace"
)

// Pathset contains the pre-computed paths used by a run of the fuzzer.
type Pathset struct {
	// DirLitmus is the directory to which litmus tests will be written.
	DirLitmus string

	// DirTrace is the directory to which traces will be written.
	DirTrace string
}

// NewPathset constructs a new pathset from the directory root.
func NewPathset(root string) *Pathset {
	return &Pathset{
		DirLitmus: path.Join(root, segLitmus),
		DirTrace:  path.Join(root, segTrace),
	}
}

// Prepare prepares this pathset by making its directories.
func (p *Pathset) Prepare() error {
	return iohelp.Mkdirs(p.DirTrace, p.DirLitmus)
}

// SubjectPaths gets the litmus and trace file paths for the subject/cycle pair c.
func (p *Pathset) SubjectPaths(c SubjectCycle) subject.FuzzFileset {
	base := c.String()
	return subject.FuzzFileset{
		Litmus: path.Join(p.DirLitmus, base+".litmus"),
		Trace:  path.Join(p.DirTrace, base+".trace"),
	}
}
