// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package normaliser_test

import (
	"path"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/model/normaliser"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
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
		olit := path.Join(root, normaliser.FileOrigLitmus)
		return testCase{
			in:  *subject.NewOrPanic(litmus.New(path.Join("foo", "bar", "baz.litmus"))),
			out: *subject.NewOrPanic(litmus.New(olit)),
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
					Litmus: *litmus.New(path.Join("barbaz", "baz.1.litmus")),
					Trace:  path.Join("barbaz", "baz.1.trace"),
				},
			},
			out: subject.Subject{
				Fuzz: &subject.Fuzz{
					Litmus: *litmus.New(r(normaliser.FileFuzzLitmus)),
					Trace:  r(normaliser.FileFuzzTrace),
				},
			},
			maps: normaliser.Map{
				r(normaliser.FileFuzzLitmus): normaliser.NewEntry(filekind.Litmus, filekind.InFuzz, "barbaz", "baz.1.litmus"),
				r(normaliser.FileFuzzTrace):  normaliser.NewEntry(filekind.Trace, filekind.InFuzz, "barbaz", "baz.1.trace"),
			},
		}
	},
	"recipe": func(root string) testCase {
		h := func(arch, file string) string { return path.Join(root, normaliser.DirRecipes, arch, file) }
		return testCase{
			in: subject.Subject{
				Recipes: map[string]recipe.Recipe{
					"arm": {
						Dir:   path.Join("burble", "armv8"),
						Files: []string{"inky.c", "pinky.c"},
					},
					"x86": {
						Dir:   path.Join("burble", "i386"),
						Files: []string{"inky.c", "pinky.c"},
					},
				},
			},
			out: subject.Subject{
				Recipes: map[string]recipe.Recipe{
					"arm": {
						Dir:   normaliser.RecipeDir(root, "arm"),
						Files: []string{"inky.c", "pinky.c"},
					},
					"x86": {
						Dir:   normaliser.RecipeDir(root, "x86"),
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
		c := func(comp, file string) string { return path.Join(root, normaliser.DirCompiles, comp, file) }
		return testCase{
			in: subject.Subject{
				Compiles: map[string]compilation.CompileResult{
					"clang": {
						Result: compilation.Result{Status: status.Ok},
						Files: compilation.CompileFileset{
							Bin: path.Join("foobaz", "clang", "a.out"),
							Log: path.Join("foobaz", "clang", "errors"),
						},
					},
					"gcc": {
						Result: compilation.Result{Status: status.Ok},
						Files: compilation.CompileFileset{
							Bin: path.Join("foobaz", "gcc", "a.out"),
							Log: path.Join("foobaz", "gcc", "errors"),
						},
					},
				},
			},
			out: subject.Subject{
				Compiles: map[string]compilation.CompileResult{
					"clang": {
						Result: compilation.Result{Status: status.Ok},
						Files: compilation.CompileFileset{
							Bin: c("clang", normaliser.FileBin),
							Log: c("clang", normaliser.FileCompileLog),
						},
					},
					"gcc": {
						Result: compilation.Result{Status: status.Ok},
						Files: compilation.CompileFileset{
							Bin: c("gcc", normaliser.FileBin),
							Log: c("gcc", normaliser.FileCompileLog),
						},
					},
				},
			},
			maps: map[string]normaliser.Entry{
				c("clang", normaliser.FileBin):        normaliser.NewEntry(filekind.Bin, filekind.InCompile, "foobaz", "clang", "a.out"),
				c("gcc", normaliser.FileBin):          normaliser.NewEntry(filekind.Bin, filekind.InCompile, "foobaz", "gcc", "a.out"),
				c("clang", normaliser.FileCompileLog): normaliser.NewEntry(filekind.Log, filekind.InCompile, "foobaz", "clang", "errors"),
				c("gcc", normaliser.FileCompileLog):   normaliser.NewEntry(filekind.Log, filekind.InCompile, "foobaz", "gcc", "errors"),
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
