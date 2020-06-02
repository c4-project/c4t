// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package recipe contains 'recipes': instructions to the machine node on how to compile a subject.
package recipe

import (
	"path"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"
)

// Recipe represents information about a lifted test recipe.
type Recipe struct {
	// Dir is the root directory of the recipe.
	Dir string `toml:"dir" json:"dir"`

	// Files is a list of files initially available in the recipe.
	Files []string `toml:"files" json:"files,omitempty"`

	// Instructions is a list of instructions for the machine node.
	Instructions []Instruction `json:"instructions,omitempty"`
}

// New constructs a recipe using the input directory dir and the options os.
func New(dir string, os ...Option) Recipe {
	r := Recipe{Dir: dir}
	Options(os...)(&r)
	return r
}

// Option is a functional option for a recipe.
type Option func(*Recipe)

// Options applies multiple options to a recipe.
func Options(os ...Option) Option {
	return func(r *Recipe) {
		for _, o := range os {
			o(r)
		}
	}
}

// AddFiles adds each file in fs to the recipe.
func AddFiles(fs ...string) Option {
	return func(r *Recipe) {
		r.Files = append(r.Files, fs...)
	}
}

// AddInstructions adds each instruction in ins to the recipe.
func AddInstructions(ins ...Instruction) Option {
	return func(r *Recipe) {
		r.Instructions = append(r.Instructions, ins...)
	}
}

// CompileFileToObj adds a set of instructions that compile the named C input to an object file.
func CompileFileToObj(file string) Option {
	return AddInstructions(
		PushInputInst(file),
		CompileObjInst(),
	)
}

// CompileAllCToExe adds a set of instructions that compile all C inputs to an executable.
func CompileAllCToExe() Option {
	return AddInstructions(
		PushInputsInst(filekind.CSrc),
		CompileExeInst(),
	)
}

// Paths retrieves the slash-joined dir/file paths for each file in the recipe.
func (r *Recipe) Paths() []string {
	paths := make([]string, len(r.Files))
	for i, f := range r.Files {
		paths[i] = path.Join(r.Dir, f)
	}
	return paths
}
