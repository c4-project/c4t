// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pathset

import (
	"github.com/c4-project/c4t/internal/helper/iohelp"
)

// Prepare prepares this pathset by making its directories.
func (p *Pathset) Prepare() error {
	return iohelp.Mkdirs(p.DirSaved, p.DirScratch)
}
