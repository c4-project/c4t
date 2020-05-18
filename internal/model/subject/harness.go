// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject

import "path"

// Harness represents information about a lifted test harness.
type Harness struct {
	// Dir is the root directory of the harness.
	Dir string `toml:"dir" json:"dir"`

	// Files is a list of files in the harness.
	Files []string `toml:"files" json:"files"`
}

// Paths retrieves the joined dir/file paths for each file in the harness.
func (h Harness) Paths() []string {
	paths := make([]string, len(h.Files))
	for i, f := range h.Files {
		paths[i] = path.Join(h.Dir, f)
	}
	return paths
}

// CPaths retrieves the joined dir/file paths for each C file in the harness.
func (h Harness) CPaths() []string {
	ps := h.Paths()
	cs := make([]string, 0, len(ps))
	for _, p := range ps {
		if path.Ext(p) == ".c" {
			cs = append(cs, p)
		}
	}
	return cs
}
