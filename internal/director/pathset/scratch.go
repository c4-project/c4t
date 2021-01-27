// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pathset

import (
	"path/filepath"

	"github.com/c4-project/c4t/internal/helper/iohelp"
)

const (
	segFuzz = "fuzz"
	segLift = "lift"
	segRun  = "run"
)

// Scratch contains the pre-computed paths for a machine run.
type Scratch struct {
	// DirFuzz is the directory to which fuzzed subjects will be output.
	DirFuzz string
	// DirLift is the directory to which lifter outputs will be written.
	DirLift string
	// DirRun is the directory into which c4t-mach output will go.
	DirRun string
}

// NewScratch creates a machine pathset rooted at root.
func NewScratch(root string) *Scratch {
	return &Scratch{
		DirFuzz: filepath.Join(root, segFuzz),
		DirLift: filepath.Join(root, segLift),
		DirRun:  filepath.Join(root, segRun),
	}
}

// Dirs gets all of the directories in this pathset, which is useful for making and removing directories.
func (p *Scratch) Dirs() []string {
	return []string{p.DirFuzz, p.DirLift, p.DirRun}
}

// Prepare prepares this pathset by making its directories.
func (p *Scratch) Prepare() error {
	return iohelp.Mkdirs(p.Dirs()...)
}
