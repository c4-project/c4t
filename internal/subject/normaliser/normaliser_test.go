// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package normaliser_test

import (
	"path"
	"testing"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/subject/normpath"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/model/litmus"

	"github.com/c4-project/c4t/internal/model/recipe"

	"github.com/c4-project/c4t/internal/model/filekind"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/subject/normaliser"

	"github.com/c4-project/c4t/internal/subject"
)

type testCase struct {
	in   subject.Subject
	out  subject.Subject
	maps normaliser.Map
}

var testSubjects = map[string]func(root string) testCase{
	"empty": func(string) testCase {
		return testCase{
			in:   subject.Subject{},
			out:  subject.Subject{},
			maps: normaliser.Map{},
		}
	},
	"litmus": func(root string) testCase {
		olit := path.Join(root, normpath.FileOrigLitmus)
		return testCase{
			in:  *subject.NewOrPanic(litmus.NewOrPanic(path.Join("foo", "bar", "baz.litmus"))),
			out: *subject.NewOrPanic(litmus.NewOrPanic(olit)),
			maps: normaliser.Map{
				olit: normaliser.NewEntry(filekind.Litmus, filekind.InOrig, "foo", "bar", "baz.litmus"),
			},
		}
	},
	"fuzz": func(root string) testCase {
		r := func(s string) string { return path.Join(root, s) }
		return testCase{
			in: subject.Subject{
				Fuzz: &subject.Fuzz{
					Litmus: *litmus.NewOrPanic(path.Join("barbaz", "baz.1.litmus")),
					Trace:  path.Join("barbaz", "baz.1.trace"),
				},
			},
			out: subject.Subject{
				Fuzz: &subject.Fuzz{
					Litmus: *litmus.NewOrPanic(r(normpath.FileFuzzLitmus)),
					Trace:  r(normpath.FileFuzzTrace),
				},
			},
			maps: normaliser.Map{
				r(normpath.FileFuzzLitmus): normaliser.NewEntry(filekind.Litmus, filekind.InFuzz, "barbaz", "baz.1.litmus"),
				r(normpath.FileFuzzTrace):  normaliser.NewEntry(filekind.Trace, filekind.InFuzz, "barbaz", "baz.1.trace"),
			},
		}
	},
	"recipe": func(root string) testCase {
		h := func(arch, file string) string { return path.Join(root, normpath.DirRecipes, arch, file) }
		return testCase{
			in: subject.Subject{
				Recipes: recipe.Map{
					id.ArchArm: {
						Dir:   path.Join("burble", "armv8"),
						Files: []string{"inky.c", "pinky.c"},
					},
					id.ArchX86: {
						Dir:   path.Join("burble", "i386"),
						Files: []string{"inky.c", "pinky.c"},
					},
				},
			},
			out: subject.Subject{
				Recipes: recipe.Map{
					id.ArchArm: {
						Dir:   normpath.RecipeDir(root, "arm"),
						Files: []string{"inky.c", "pinky.c"},
					},
					id.ArchX86: {
						Dir:   normpath.RecipeDir(root, "x86"),
						Files: []string{"inky.c", "pinky.c"},
					},
				},
			},
			maps: normaliser.Map{
				h("arm", "inky.c"):  normaliser.NewEntry(filekind.CSrc, filekind.InRecipe, "burble", "armv8", "inky.c"),
				h("arm", "pinky.c"): normaliser.NewEntry(filekind.CSrc, filekind.InRecipe, "burble", "armv8", "pinky.c"),
				h("x86", "inky.c"):  normaliser.NewEntry(filekind.CSrc, filekind.InRecipe, "burble", "i386", "inky.c"),
				h("x86", "pinky.c"): normaliser.NewEntry(filekind.CSrc, filekind.InRecipe, "burble", "i386", "pinky.c"),
			},
		}
	},
	"compile": func(root string) testCase {
		c := func(comp, file string) string { return path.Join(root, normpath.DirCompiles, comp, file) }
		return testCase{
			in: subject.Subject{
				Compilations: compilation.Map{
					id.FromString("clang"): {
						Compile: &compilation.CompileResult{
							Result: compilation.Result{Status: status.Ok},
							Files: compilation.CompileFileset{
								Bin: path.Join("foobaz", "clang", "a.out"),
								Log: path.Join("foobaz", "clang", "errors"),
							},
						},
					},
					id.FromString("gcc"): {
						Compile: &compilation.CompileResult{
							Result: compilation.Result{Status: status.Ok},
							Files: compilation.CompileFileset{
								Bin: path.Join("foobaz", "gcc", "a.out"),
								Log: path.Join("foobaz", "gcc", "errors"),
							},
						},
					},
				},
			},
			out: subject.Subject{
				Compilations: compilation.Map{
					id.FromString("clang"): {
						Compile: &compilation.CompileResult{
							Result: compilation.Result{Status: status.Ok},
							Files: compilation.CompileFileset{
								Bin: c("clang", normpath.FileBin),
								Log: c("clang", normpath.FileCompileLog),
							},
						},
					},
					id.FromString("gcc"): {
						Compile: &compilation.CompileResult{
							Result: compilation.Result{Status: status.Ok},
							Files: compilation.CompileFileset{
								Bin: c("gcc", normpath.FileBin),
								Log: c("gcc", normpath.FileCompileLog),
							},
						},
					},
				},
			},
			maps: map[string]normaliser.Entry{
				c("clang", normpath.FileBin):        normaliser.NewEntry(filekind.Bin, filekind.InCompile, "foobaz", "clang", "a.out"),
				c("gcc", normpath.FileBin):          normaliser.NewEntry(filekind.Bin, filekind.InCompile, "foobaz", "gcc", "a.out"),
				c("clang", normpath.FileCompileLog): normaliser.NewEntry(filekind.Log, filekind.InCompile, "foobaz", "clang", "errors"),
				c("gcc", normpath.FileCompileLog):   normaliser.NewEntry(filekind.Log, filekind.InCompile, "foobaz", "gcc", "errors"),
			},
		}
	},
}

// TestNormaliser_Normalise checks the normaliser on various small subject cases.
func TestNormaliser_Normalise(t *testing.T) {
	t.Parallel()

	for name, c := range testSubjects {
		c := c("root")
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			n := normaliser.New("root")
			s, err := n.Normalise(c.in)
			if assert.NoError(t, err) {
				assert.Equal(t, c.out, *s)
				assert.Equal(t, c.maps, n.Mappings)
			}
		})
	}
}
