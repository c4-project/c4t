// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package normaliser

import "path"

const (
	// FileBin is the normalised name for output binaries.
	FileBin = "a.out"
	// FileCompileLog is the normalised name for compilation logs.
	FileCompileLog = "compile.log"
	// FileOrigLitmus is the normalised name for pre-fuzz litmus tests.
	FileOrigLitmus = "orig.litmus"
	// FileFuzzLitmus is the normalised name for post-fuzz litmus tests.
	FileFuzzLitmus = "fuzz.litmus"
	// FileFuzzTrace is the normalised name for fuzzer traces.
	FileFuzzTrace = "fuzz.trace"
	// DirCompiles is the normalised directory for compile results.
	DirCompiles = "compiles"
	// DirRecipes is the normalised directory for recipe results.
	DirRecipes = "recipes"
)

// RecipeDir gets the normalised recipe directory under root and for architecture ID-string arch.
func RecipeDir(root, arch string) string {
	return path.Join(root, DirRecipes, arch)
}
