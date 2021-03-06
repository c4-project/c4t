// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pathset

import "github.com/c4-project/c4t/internal/stage/analyser/saver"

// Instance is an instance-specific pathset.
type Instance struct {
	// SavedPaths contains the save pathset for this machine.
	Saved saver.Pathset
	// ScratchPaths contains the scratch pathset for this machine.
	Scratch Scratch
}
