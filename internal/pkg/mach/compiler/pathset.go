// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"path"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
)

const (
	segBins = "bins"
	segLogs = "logs"
)

// Pathset contains the various directories used by the test compiler.
type Pathset struct {
	// DirBins is the directory into which compiled binaries should go.
	DirBins string

	// DirLogs is the directory into which compiler logs should go.
	DirLogs string
}

// NewPathset constructs a new pathset from the directory root.
func NewPathset(root string) *Pathset {
	return &Pathset{
		DirBins: path.Join(root, segBins),
		DirLogs: path.Join(root, segLogs),
	}
}

// Prepare prepares this pathset by making its directories.
// It takes a slice of compilers for which directories should be made.
func (p *Pathset) Prepare(compilers []id.ID) error {
	return iohelp.Mkdirs(p.Dirs(compilers...)...)
}

// Dirs gets all of the directories involved in a pathset over compiler ID set compilers.
func (p *Pathset) Dirs(compilers ...id.ID) []string {
	roots := []string{p.DirBins, p.DirLogs}
	dirs := make([]string, 0, (len(compilers)+1)*len(roots))
	for _, root := range roots {
		dirs = append(dirs, root)
		for _, c := range compilers {
			elems := append([]string{root}, c.Tags()...)
			dirs = append(dirs, path.Join(elems...))
		}
	}
	return dirs
}

// SubjectPaths gets the binary and log file paths for the subject/compiler pair sc.
func (p *Pathset) SubjectPaths(sc SubjectCompile) subject.CompileFileset {
	csub := append(sc.CompilerID.Tags(), sc.Name)
	bpath := append([]string{p.DirBins}, csub...)
	lpath := append([]string{p.DirLogs}, csub...)
	return subject.CompileFileset{Bin: path.Join(bpath...), Log: path.Join(lpath...)}
}
