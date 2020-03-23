// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pathset

import "github.com/MattWindsor91/act-tester/internal/pkg/helpers/iohelp"

// All of these methods alter the filesystem, and so are quite hard to test...!

// Prepare prepares this pathset by making its directories.
func (p *Pathset) Prepare() error {
	return iohelp.Mkdirs(p.DirSaved, p.DirScratch)
}

// Prepare prepares this pathset by making its directories.
func (p *Scratch) Prepare() error {
	return iohelp.Mkdirs(p.DirPlan, p.DirFuzz, p.DirLift, p.DirRun)
}

// Prepare prepares this pathset by making its directories.
func (s *Saved) Prepare() error {
	return iohelp.Mkdirs(s.DirCompileFailures, s.DirFlagged, s.DirRunFailures, s.DirTimeouts)
}
