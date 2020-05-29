// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
)

// ExampleSubject_BestLitmus is a testable example for BestLitmus.
func ExampleSubject_BestLitmus() {
	s1 := subject.Subject{OrigLitmus: "foo.litmus"}
	b1, _ := s1.BestLitmus()

	// This subject has a fuzzed litmus file, which takes priority.
	s2 := subject.Subject{OrigLitmus: "foo.litmus", Fuzz: &subject.Fuzz{Files: subject.FuzzFileset{Litmus: "bar.litmus"}}}
	b2, _ := s2.BestLitmus()

	fmt.Println("s1:", b1)
	fmt.Println("s2:", b2)

	// Output:
	// s1: foo.litmus
	// s2: bar.litmus
}

// ExampleSubject_CompileResult is a testable example for CompileResult.
func ExampleSubject_CompileResult() {
	s := subject.Subject{Compiles: map[string]subject.CompileResult{
		"gcc":   {Result: subject.Result{Status: status.Ok}, Files: subject.CompileFileset{Bin: "a.out", Log: "gcc.log"}},
		"clang": {Result: subject.Result{Status: status.CompileFail}, Files: subject.CompileFileset{Bin: "a.out", Log: "clang.log"}},
	}}
	gr, _ := s.CompileResult(id.FromString("gcc"))
	cr, _ := s.CompileResult(id.FromString("clang"))

	fmt.Println("gcc:", gr.Status, gr.Files.Bin, gr.Files.Log)
	fmt.Println("clang:", cr.Status, cr.Files.Bin, cr.Files.Log)

	// Output:
	// gcc: Ok a.out gcc.log
	// clang: CompileFail a.out clang.log
}

// ExampleSubject_Harness is a testable example for Recipe.
func ExampleSubject_Harness() {
	s := subject.Subject{Harnesses: map[string]recipe.Recipe{
		"x86.64": {Dir: "foo", Files: []string{"bar", "baz"}},
		"arm":    {Dir: "foobar", Files: []string{"barbaz"}},
	}}
	xs, _ := s.Harness(id.ArchX8664)
	as, _ := s.Harness(id.ArchArm)

	for _, r := range xs.Files {
		fmt.Println(r)
	}
	for _, r := range as.Files {
		fmt.Println(r)
	}

	// Output:
	// bar
	// baz
	// barbaz
}

// ExampleSubject_RunOf is a testable example for RunOf.
func ExampleSubject_RunOf() {
	s := subject.Subject{Runs: map[string]subject.RunResult{
		"gcc":   {Result: subject.Result{Status: status.Ok}},
		"clang": {Result: subject.Result{Status: status.RunTimeout}},
	}}
	gr, _ := s.RunOf(id.FromString("gcc"))
	cr, _ := s.RunOf(id.FromString("clang"))

	fmt.Println("gcc:", gr.Status)
	fmt.Println("clang:", cr.Status)

	// Output:
	// gcc: Ok
	// clang: RunTimeout
}

// TestSubject_CompileResult_Missing checks that trying to get a compile for a missing compiler triggers
// the appropriate error.
func TestSubject_CompileResult_Missing(t *testing.T) {
	var s subject.Subject
	_, err := s.CompileResult(id.FromString("gcc"))
	testhelp.ExpectErrorIs(t, err, subject.ErrMissingCompile, "missing compile result path")
}

// TestSubject_AddCompileResult checks that AddCompileResult is working properly.
func TestSubject_AddCompileResult(t *testing.T) {
	var s subject.Subject
	c := subject.CompileResult{
		Result: subject.Result{Status: status.Ok},
		Files: subject.CompileFileset{
			Bin: "a.out",
			Log: "gcc.log",
		},
	}

	mcomp := id.FromString("gcc")

	t.Run("initial-add", func(t *testing.T) {
		if err := s.AddCompileResult(mcomp, c); err != nil {
			t.Fatalf("err when adding compile to empty subject: %v", err)
		}
	})
	t.Run("add-get", func(t *testing.T) {
		c2, err := s.CompileResult(mcomp)
		if err != nil {
			t.Fatalf("err when getting added compile: %v", err)
		}
		if !reflect.DeepEqual(c2, c) {
			t.Fatalf("added compile (%v) came back wrong (%v)", c2, c)
		}
	})
	t.Run("add-dupe", func(t *testing.T) {
		err := s.AddCompileResult(mcomp, subject.CompileResult{})
		testhelp.ExpectErrorIs(t, err, subject.ErrDuplicateCompile, "adding compile twice")
	})
}

// TestSubject_Harness_Missing checks that trying to get a harness path for a missing arch triggers
// the appropriate error.
func TestSubject_Harness_Missing(t *testing.T) {
	var s subject.Subject
	_, err := s.Harness(id.FromString("x86.64"))
	testhelp.ExpectErrorIs(t, err, subject.ErrMissingHarness, "missing harness path")
}

// TestSubject_AddHarness checks that AddHarness is working properly.
func TestSubject_AddHarness(t *testing.T) {
	var s subject.Subject
	h := recipe.Recipe{
		Dir:   "foo",
		Files: []string{"bar", "baz"},
	}

	march := id.ArchX8664

	t.Run("initial-add", func(t *testing.T) {
		if err := s.AddHarness(march, h); err != nil {
			t.Fatalf("err when adding harness to empty subject: %v", err)
		}
	})
	t.Run("add-get", func(t *testing.T) {
		h2, err := s.Harness(march)
		if err != nil {
			t.Fatalf("err when getting added harness: %v", err)
		}
		if !reflect.DeepEqual(h2, h) {
			t.Fatalf("added harness (%v) came back wrong (%v)", h2, h)
		}
	})
	t.Run("add-dupe", func(t *testing.T) {
		err := s.AddHarness(march, recipe.Recipe{})
		testhelp.ExpectErrorIs(t, err, subject.ErrDuplicateHarness, "adding harness twice")
	})
}

// TestSubject_RunOf_Missing checks that trying to get a run for a missing compiler gives
// the appropriate error.
func TestSubject_RunOf_Missing(t *testing.T) {
	var s subject.Subject
	_, err := s.RunOf(id.FromString("gcc"))
	testhelp.ExpectErrorIs(t, err, subject.ErrMissingRun, "missing run result path")
}

// TestSubject_AddRun checks that AddRun is working properly.
func TestSubject_AddRun(t *testing.T) {
	var s subject.Subject
	c := subject.RunResult{Result: subject.Result{Status: status.RunTimeout}}

	mcomp := id.FromString("gcc")

	t.Run("initial-add", func(t *testing.T) {
		if err := s.AddRun(mcomp, c); err != nil {
			t.Fatalf("err when adding run to empty subject: %v", err)
		}
	})
	t.Run("add-get", func(t *testing.T) {
		c2, err := s.RunOf(mcomp)
		if err != nil {
			t.Fatalf("err when getting added run: %v", err)
		}
		if !reflect.DeepEqual(c2, c) {
			t.Fatalf("added run (%v) came back wrong (%v)", c2, c)
		}
	})
	t.Run("add-dupe", func(t *testing.T) {
		err := s.AddRun(mcomp, subject.RunResult{})
		testhelp.ExpectErrorIs(t, err, subject.ErrDuplicateRun, "adding compile twice")
	})
}

// TestSubject_BestLitmus tests a few cases of BestLitmus.
// It should be more comprehensive than the examples.
func TestSubject_BestLitmus(t *testing.T) {
	// Note that the presence of 'err' overrides 'want'.
	cases := map[string]struct {
		s    subject.Subject
		err  error
		want string
	}{
		"zero":             {s: subject.Subject{}, err: subject.ErrNoBestLitmus, want: ""},
		"zero-fuzz":        {s: subject.Subject{Fuzz: &subject.Fuzz{}}, err: subject.ErrNoBestLitmus, want: ""},
		"litmus-only":      {s: subject.Subject{OrigLitmus: "foo"}, err: nil, want: "foo"},
		"litmus-only-fuzz": {s: subject.Subject{OrigLitmus: "foo", Fuzz: &subject.Fuzz{}}, err: nil, want: "foo"},
		"fuzz":             {s: subject.Subject{OrigLitmus: "foo", Fuzz: &subject.Fuzz{Files: subject.FuzzFileset{Litmus: "bar"}}}, err: nil, want: "bar"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := c.s.BestLitmus()
			switch {
			case err != nil && c.err == nil:
				t.Errorf("unexpected BestLitmus(%v) error: %v", c.s, err)
			case err != nil && !errors.Is(err, c.err):
				t.Errorf("wrong BestLitmus(%v) error: got %v; want %v", c.s, err, c.err)
			case err == nil && c.err != nil:
				t.Errorf("no BestLitmus(%v) error; want %v", c.s, err)
			case err == nil && got != c.want:
				t.Errorf("BestLitmus(%v)=%q; want %q", c.s, got, c.want)
			}
		})
	}
}
