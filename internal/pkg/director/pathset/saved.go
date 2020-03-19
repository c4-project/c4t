// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pathset

import (
	"path/filepath"
)

const (
	segCompileFailures = "compile_fail"
	segFlagged         = "flagged"
	segRunFailures     = "run_fail"
	segTimeouts        = "timeout"
)

// Saved contains the pre-computed paths for saving 'interesting' run results.
type Saved struct {
	// DirCompileFailures stores subjects that ran into compiler failures.
	DirCompileFailures string
	// DirFlagged stores subjects that have been flagged as possibly buggy.
	DirFlagged string
	// DirRunFailures stores subjects that met run failures.
	DirRunFailures string
	// DirTimeouts stores subjects that timed out.
	DirTimeouts string
}

// NewSaved creates a save pathset rooted at root.
func NewSaved(root string) *Saved {
	return &Saved{
		DirCompileFailures: filepath.Join(root, segCompileFailures),
		DirFlagged:         filepath.Join(root, segFlagged),
		DirRunFailures:     filepath.Join(root, segRunFailures),
		DirTimeouts:        filepath.Join(root, segTimeouts),
	}
}
