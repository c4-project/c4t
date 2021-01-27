// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package filekind

import (
	"path/filepath"
	"strings"
)

// GuessFromFile tries to guess the kind of the file whose filepath is path.
// It does so only from the file extension.
func GuessFromFile(fpath string) Kind {
	if fpath == "a.out" {
		return Bin
	}

	switch strings.TrimLeft(filepath.Ext(fpath), ".") {
	case "c":
		return CSrc
	case "h":
		return CHeader
	case "litmus":
		return Litmus
	case "trace":
		return Trace
	}

	return Other
}

// FilterFiles filters fpaths to paths whose guess matches this kind.
func (k Kind) FilterFiles(fpaths []string) []string {
	out := make([]string, 0, len(fpaths))
	for _, f := range fpaths {
		if GuessFromFile(f).Matches(k) {
			out = append(out, f)
		}
	}
	return out
}
