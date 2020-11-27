// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import "path"

const (
	segObs      = "obs"
	segFailures = "failures"
)

// Pathset contains the top-level paths for a runner instance.
type Pathset struct {
	// DirObs is the directory into which observations should go.
	DirObs string

	// DirFailures is the directory into which failing tests should go.
	DirFailures string
}

// NewPathset constructs a new pathset from the directory root.
func NewPathset(root string) *Pathset {
	return &Pathset{
		DirObs:      path.Join(root, segObs),
		DirFailures: path.Join(root, segFailures),
	}
}
