// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package pathset contains the various path-sets for the director.
// These spill into a separate package on the basis of there being so many of them.
package pathset

import (
	"path/filepath"

	"github.com/c4-project/c4t/internal/stage/analyser/saver"

	"github.com/c4-project/c4t/internal/model/id"
)

const (
	segSaved   = "saved"
	segScratch = "scratch"
)

// Pathset contains the pre-computed paths used by the director.
type Pathset struct {
	// DirSaved is the directory into which saved runs get copied.
	DirSaved string

	// DirScratch is the directory that the director uses for ephemeral run data.
	DirScratch string
}

// New constructs a new pathset from the directory root.
func New(root string) *Pathset {
	return &Pathset{
		DirSaved:   filepath.Join(root, segSaved),
		DirScratch: filepath.Join(root, segScratch),
	}
}

// Instance gets the instance pathset for a machine with ID mid.
func (p Pathset) Instance(mid id.ID) *Instance {
	tags := mid.Tags()
	saved := append([]string{p.DirSaved}, tags...)
	scratch := append([]string{p.DirScratch}, tags...)
	// TODO(@MattWindsor91): the pointer soup here needs simplifying
	return &Instance{
		Saved:   *saver.NewPathset(filepath.Join(saved...)),
		Scratch: *NewScratch(filepath.Join(scratch...)),
	}
}
