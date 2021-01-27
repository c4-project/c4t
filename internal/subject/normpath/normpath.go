// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package normpath contains various 'normalised path' fragments, and functions for constructing them.
// These are used by both the normaliser (which rewrites a plan to use normalised paths) and various other parts of
// the tooling (which assumes the use of normalised paths).
package normpath

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

	// TarSuffix is the extension used by the saver when saving tarballs, and presumed by archive-transparent file readers.
	TarSuffix = ".tar.gz"
)

// RecipeDir gets the normalised recipe directory under root and for architecture ID-string arch.
func RecipeDir(root, arch string) string {
	return path.Join(root, DirRecipes, arch)
}
