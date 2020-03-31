// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package pathset contains the various path-sets for the director.
// These spill into a separate package on the basis of there being so many of them.
package pathset

import (
	"path/filepath"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

const (
	segFuzz    = "fuzz"
	segLift    = "lift"
	segPlan    = "plan"
	segRun     = "run"
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

// MachineSave gets the saved-subjects pathset for a machine with ID mid.
func (p *Pathset) MachineSaved(mid id.ID) *Saved {
	segs := append([]string{p.DirSaved}, mid.Tags()...)
	return NewSaved(filepath.Join(segs...))
}

// MachineScratch gets the scratch pathset for a machine with ID mid.
func (p *Pathset) MachineScratch(mid id.ID) *Scratch {
	segs := append([]string{p.DirScratch}, mid.Tags()...)
	return NewScratch(filepath.Join(segs...))
}
