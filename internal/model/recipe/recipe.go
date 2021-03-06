// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package recipe contains 'recipes': instructions to the machine node on how to compile a subject.
package recipe

import (
	"path"

	"github.com/c4-project/c4t/internal/id"
)

// Recipe represents information about a lifted test recipe.
type Recipe struct {
	// Dir is the root directory of the recipe.
	Dir string `toml:"dir" json:"dir"`

	// Files is a list of files initially available in the recipe.
	Files []string `toml:"files" json:"files,omitempty"`

	// Instructions is a list of instructions for the machine node.
	Instructions []Instruction `json:"instructions,omitempty"`

	// OutType is the type of output this recipe promises.  The output file is implicit.
	Output Output `json:"out_type,omitempty"`
}

// New constructs a recipe using the input directory dir and the options os.
func New(dir string, otype Output, os ...Option) (Recipe, error) {
	r := Recipe{Dir: dir, Output: otype}
	if err := Options(os...)(&r); err != nil {
		return Recipe{}, err
	}
	return r, nil
}

// Paths retrieves the slash-joined dir/file paths for each file in the recipe.
func (r *Recipe) Paths() []string {
	paths := make([]string, len(r.Files))
	for i, f := range r.Files {
		paths[i] = path.Join(r.Dir, f)
	}
	return paths
}

// NeedsCompile gets whether this recipe needs to be compiled (ie, its instructions should be interpreted).
func (r *Recipe) NeedsCompile() bool {
	return r.Output != OutNothing
}

// Map is shorthand from a map from architecture IDs to recipes.
type Map map[id.ID]Recipe
