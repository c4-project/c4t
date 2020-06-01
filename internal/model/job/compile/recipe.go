// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compile

import (
	"github.com/MattWindsor91/act-tester/internal/model/compiler"
	"github.com/MattWindsor91/act-tester/internal/model/filekind"
	"github.com/MattWindsor91/act-tester/internal/model/recipe"
)

// Recipe represents a request to compile a multi-stage recipe with a particular compiler.
type Recipe struct {
	Compile

	// Instructions contains the instruction set for the compiler job, if the job kind is CompileRecipe.
	Instructions []recipe.Instruction
}

// FromRecipe constructs a recipe compile from the recipe r, compiler c, and output file out.
func FromRecipe(c *compiler.Compiler, r recipe.Recipe, out string) Recipe {
	return Recipe{
		Compile: Compile{
			Compiler: c,
			In:       filekind.CSrc.FilterFiles(r.Paths()),
			Out:      out,
		},
		Instructions: r.Instructions,
	}
}
