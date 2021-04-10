// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package subject_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/model/litmus"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/model/recipe"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/subject"

	"github.com/c4-project/c4t/internal/helper/testhelp"
)

// ExampleSubject_BestLitmus is a testable example for BestLitmus.
func ExampleSubject_BestLitmus() {
	l, _ := litmus.New("foo.litmus")

	s1, _ := subject.New(l)
	b1, _ := s1.BestLitmus()

	// This subject has a fuzzed litmus file, which takes priority.
	f, _ := litmus.New("bar.litmus")
	s2, _ := subject.New(l, subject.WithFuzz(&subject.Fuzz{Litmus: *f}))
	b2, _ := s2.BestLitmus()

	fmt.Println("s1:", b1.Path)
	fmt.Println("s2:", b2.Path)

	// Output:
	// s1: foo.litmus
	// s2: bar.litmus
}

// ExampleSubject_CompileResult is a testable example for CompileResult.
func ExampleSubject_CompileResult() {
	s := subject.Subject{Compilations: compilation.Map{
		id.FromString("gcc"): {
			Compile: &compilation.CompileResult{
				Result: compilation.Result{Status: status.Ok}, Files: compilation.CompileFileset{Bin: "a.out", Log: "gcc.log"},
			}},
		id.FromString("clang"): {
			Compile: &compilation.CompileResult{
				Result: compilation.Result{Status: status.CompileFail}, Files: compilation.CompileFileset{Bin: "a.out", Log: "clang.log"},
			}},
	}}
	gr, _ := s.CompileResult(id.FromString("gcc"))
	cr, _ := s.CompileResult(id.FromString("clang"))

	fmt.Println("gcc:", gr.Status, gr.Files.Bin, gr.Files.Log)
	fmt.Println("clang:", cr.Status, cr.Files.Bin, cr.Files.Log)

	// Output:
	// gcc: Ok a.out gcc.log
	// clang: CompileFail a.out clang.log
}

// ExampleSubject_Recipe is a testable example for Recipe.
func ExampleSubject_Recipe() {
	s := subject.Subject{Recipes: recipe.Map{
		id.ArchX8664: {Dir: "foo", Files: []string{"bar", "baz"}},
		id.ArchArm:   {Dir: "foobar", Files: []string{"barbaz"}},
	}}
	xsn, xs, _ := s.Recipe(id.ArchX8664)
	asn, as, _ := s.Recipe(id.ArchArm)

	fmt.Println("#", xsn)
	for _, r := range xs.Files {
		fmt.Println(r)
	}
	fmt.Println("#", asn)
	for _, r := range as.Files {
		fmt.Println(r)
	}

	// Output:
	// # x86.64
	// bar
	// baz
	// # arm
	// barbaz
}

// ExampleSubject_RunResult is a testable example for Subject.RunResult.
func ExampleSubject_RunResult() {
	s := subject.Subject{Compilations: compilation.Map{
		id.FromString("gcc"): {
			Run: &compilation.RunResult{
				Result: compilation.Result{Status: status.Ok},
			}},
		id.FromString("clang"): {
			Run: &compilation.RunResult{
				Result: compilation.Result{Status: status.RunTimeout},
			}},
	}}
	gr, _ := s.RunResult(id.FromString("gcc"))
	cr, _ := s.RunResult(id.FromString("clang"))

	fmt.Println("gcc:", gr.Status)
	fmt.Println("clang:", cr.Status)

	// Output:
	// gcc: Ok
	// clang: RunTimeout
}

// TestSubject_CompileResult_Missing checks that trying to get a compile for a missing compiler triggers
// the appropriate errors.
func TestSubject_CompileResult_Missing(t *testing.T) {
	t.Parallel()

	var s subject.Subject
	_, err := s.CompileResult(id.FromString("gcc"))
	testhelp.ExpectErrorIs(t, err, subject.ErrMissingCompilation, "missing compilations")

	s.Compilations = compilation.Map{id.FromString("gcc"): {}}
	_, err = s.CompileResult(id.FromString("gcc"))
	testhelp.ExpectErrorIs(t, err, subject.ErrMissingCompile, "missing compile result path")
}

// TestSubject_AddCompileResult checks that AddCompileResult is working properly.
func TestSubject_AddCompileResult(t *testing.T) {
	t.Parallel()

	var s subject.Subject
	c := compilation.CompileResult{
		Result:   compilation.Result{Status: status.Ok},
		RecipeID: id.ArchArm,
		Files: compilation.CompileFileset{
			Bin: "a.out",
			Log: "gcc.log",
		},
	}

	mcomp := id.FromString("gcc")

	t.Run("initial-add", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.AddCompileResult(mcomp, c), "err when adding compile to empty subject")
	})
	t.Run("add-get", func(t *testing.T) {
		t.Parallel()
		c2, err := s.CompileResult(mcomp)
		if assert.NoError(t, err, "err when getting added compile") {
			assert.Equalf(t, c, *c2, "added compile (%v) came back wrong (%v)", c2, c)
		}
	})
	t.Run("add-dupe", func(t *testing.T) {
		t.Parallel()
		err := s.AddCompileResult(mcomp, compilation.CompileResult{})
		testhelp.ExpectErrorIs(t, err, subject.ErrDuplicateCompile, "adding compile twice")
	})
}

// TestSubject_Recipe_missing checks that trying to get a recipe for a missing arch triggers
// the appropriate error.
func TestSubject_Recipe_missing(t *testing.T) {
	var s subject.Subject
	_, _, err := s.Recipe(id.FromString("x86.64"))
	testhelp.ExpectErrorIs(t, err, subject.ErrMissingRecipe, "missing recipe path")
}

// TestSubject_AddRecipe checks that AddRecipe is working properly.
func TestSubject_AddRecipe(t *testing.T) {
	t.Parallel()

	var s subject.Subject
	h := recipe.Recipe{
		Dir:   "foo",
		Files: []string{"bar", "baz"},
	}

	march := id.ArchX8664

	t.Run("initial-add", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.AddRecipe(march, h), "err when adding recipe to empty subject")
	})
	t.Run("add-get", func(t *testing.T) {
		t.Parallel()
		m2, h2, err := s.Recipe(march)
		if assert.NoError(t, err, "err when getting added recipe") {
			assert.Equal(t, march, m2, "wrong recipe ID")
			assert.Equalf(t, h, h2, "added recipe (%v) came back wrong (%v)", h2, h)
		}
	})
	t.Run("add-dupe", func(t *testing.T) {
		t.Parallel()
		err := s.AddRecipe(march, recipe.Recipe{})
		testhelp.ExpectErrorIs(t, err, subject.ErrDuplicateRecipe, "adding recipe twice")
	})
}

// TestSubject_RunOf_Missing checks that trying to get a run for a missing compiler gives
// the appropriate error.
func TestSubject_RunOf_Missing(t *testing.T) {
	t.Parallel()

	var s subject.Subject
	_, err := s.RunResult(id.FromString("gcc"))
	testhelp.ExpectErrorIs(t, err, subject.ErrMissingCompilation, "missing compilation")

	s.Compilations = compilation.Map{id.FromString("gcc"): {}}
	_, err = s.RunResult(id.FromString("gcc"))
	testhelp.ExpectErrorIs(t, err, subject.ErrMissingRun, "missing run result")
}

// TestSubject_AddRun checks that AddRun is working properly.
func TestSubject_AddRun(t *testing.T) {
	t.Parallel()

	var s subject.Subject
	c := compilation.RunResult{Result: compilation.Result{Status: status.RunTimeout}}

	mcomp := id.FromString("gcc")

	t.Run("initial-add", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.AddRun(mcomp, c), "err when adding run to empty subject")
	})
	t.Run("add-get", func(t *testing.T) {
		t.Parallel()
		c2, err := s.RunResult(mcomp)
		if assert.NoError(t, err, "err when getting added run") {
			assert.Equalf(t, c, *c2, "added run (%v) came back wrong (%v)", c2, c)
		}
	})
	t.Run("add-dupe", func(t *testing.T) {
		t.Parallel()
		err := s.AddRun(mcomp, compilation.RunResult{})
		testhelp.ExpectErrorIs(t, err, subject.ErrDuplicateRun, "adding compile twice")
	})
}

// TestSubject_BestLitmus tests a few cases of BestLitmus.
// It should be more comprehensive than the examples.
func TestSubject_BestLitmus(t *testing.T) {
	t.Parallel()

	// Note that the presence of 'err' overrides 'want'.
	cases := map[string]struct {
		s    subject.Subject
		err  error
		want string
	}{
		"zero":             {s: subject.Subject{}, err: subject.ErrNoBestLitmus, want: ""},
		"zero-fuzz":        {s: subject.Subject{Fuzz: &subject.Fuzz{}}, err: subject.ErrNoBestLitmus, want: ""},
		"litmus-only":      {s: *subject.NewOrPanic(litmus.NewOrPanic("foo")), err: nil, want: "foo"},
		"litmus-only-fuzz": {s: *subject.NewOrPanic(litmus.NewOrPanic("foo"), subject.WithFuzz(&subject.Fuzz{})), err: nil, want: "foo"},
		"fuzz":             {s: *subject.NewOrPanic(litmus.NewOrPanic("foo"), subject.WithFuzz(&subject.Fuzz{Litmus: *litmus.NewOrPanic("bar")})), err: nil, want: "bar"},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := c.s.BestLitmus()
			switch {
			case err != nil && c.err == nil:
				t.Errorf("unexpected BestLitmus(%v) error: %v", c.s, err)
			case err != nil && !errors.Is(err, c.err):
				t.Errorf("wrong BestLitmus(%v) error: got %v; want %v", c.s, err, c.err)
			case err == nil && c.err != nil:
				t.Errorf("no BestLitmus(%v) error; want %v", c.s, err)
			case err == nil && got.Path != c.want:
				t.Errorf("BestLitmus(%v)=%q; want %q", c.s, got.Path, c.want)
			}
		})
	}
}
