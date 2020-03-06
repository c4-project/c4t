// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"path"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
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

// NewPathset constructs a new pathset from the directory root.
func NewPathset(root string) *Pathset {
	return &Pathset{
		DirSaved:   path.Join(root, segSaved),
		DirScratch: path.Join(root, segScratch),
	}
}

// Prepare prepares this pathset by making its directories.
func (p *Pathset) Prepare() error {
	return iohelp.Mkdirs(p.DirSaved, p.DirScratch)
}
