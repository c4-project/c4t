// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package recipe contains 'recipes': instructions to the machine node on how to compile a subject.
package recipe

import "path"

// Recipe represents information about a lifted test recipe.
type Recipe struct {
	// Dir is the root directory of the recipe.
	Dir string `toml:"dir" json:"dir"`

	// Files is a list of files initially available in the recipe.
	Files []string `toml:"files" json:"files,omitempty"`

	// Instructions is a list of instructions for the machine node.
	Instructions []Instruction `json:"instructions,omitempty"`
}

// Paths retrieves the joined dir/file paths for each file in the recipe.
func (r *Recipe) Paths() []string {
	paths := make([]string, len(r.Files))
	for i, f := range r.Files {
		paths[i] = path.Join(r.Dir, f)
	}
	return paths
}
