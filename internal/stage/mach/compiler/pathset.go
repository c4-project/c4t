// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"path"
	"path/filepath"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/helper/iohelp"
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
		DirBins: filepath.Join(root, segBins),
		DirLogs: filepath.Join(root, segLogs),
	}
}

// Prepare prepares this pathset for compilers cs by making its directories.
// It takes compilers for which directories should be made.
func (p *Pathset) Prepare(cs ...id.ID) error {
	// TODO(@MattWindsor91): make a record of the directories, and error if we try to use different ones in SubjectPaths.
	return iohelp.Mkdirs(p.Dirs(cs...)...)
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

// SubjectPaths gets the compilation fileset for the subject/compiler pair sc.
func (p *Pathset) SubjectPaths(sc compilation.Name) compilation.CompileFileset {
	csub := sc.Path()
	bpath := append([]string{filepath.ToSlash(p.DirBins)}, csub)
	lpath := append([]string{filepath.ToSlash(p.DirLogs)}, csub)
	return compilation.CompileFileset{Bin: path.Join(bpath...), Log: path.Join(lpath...)}
}
