// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pathset

import "github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"

// Instance is an instance-specific pathset.
type Instance struct {
	// SavedPaths contains the save pathset for this machine.
	Saved saver.Pathset
	// ScratchPaths contains the scratch pathset for this machine.
	Scratch Scratch
}
