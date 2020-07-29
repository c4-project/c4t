// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus

import (
	"path"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Mock produces a representative corpus including the following features:
// - a subject with a failed compilation;
// - a subject with a flagged observation.
func Mock() Corpus {
	return Corpus{
		"foo":    *subject.NewOrPanic(litmus.New("foo.litmus", litmus.WithThreads(1))),
		"bar":    *MockFailedCompile("bar"),
		"baz":    *MockFlaggedRun("baz"),
		"barbaz": *MockTimeoutRun("barbaz"),
	}
}

// MockFailedCompile expands to a realistic looking subject that contains a failed compilation.
func MockFailedCompile(name string) *subject.Subject {
	return subject.NewOrPanic(
		litmus.New(name+".litmus", litmus.WithThreads(8)),
		subject.WithRecipe(id.ArchArm,
			recipe.Recipe{
				Dir:   "arm",
				Files: []string{"run.c", "aux.c", "aux.h"},
			},
		),
		subject.WithCompile(id.FromString("gcc"),
			compilation.CompileResult{
				Result: compilation.Result{Status: status.CompileFail},
				Files:  compilation.CompileFileset{},
			},
		),
		subject.WithCompile(id.FromString("clang"),
			MockSuccessfulCompile("clang", name),
		),
		subject.WithRun(id.FromString("gcc"),
			compilation.RunResult{Result: compilation.Result{Status: status.CompileFail}},
		),
		subject.WithRun(id.FromString("clang"),
			compilation.RunResult{Result: compilation.Result{Status: status.Ok}},
		),
	)
}

// MockFlaggedRun expands to a realistic looking subject that contains some flagged runs.
func MockFlaggedRun(name string) *subject.Subject {
	return subject.NewOrPanic(
		litmus.New(name+".litmus", litmus.WithThreads(2)),
		subject.WithRecipe(id.ArchX8664, MockRecipe("x86")),
		subject.WithCompile(id.FromString("gcc"), MockSuccessfulCompile("gcc", name)),
		subject.WithCompile(id.FromString("icc"), MockSuccessfulCompile("icc", name)),
		subject.WithRun(id.FromString("gcc"), compilation.RunResult{Result: compilation.Result{Status: status.Flagged}}),
		subject.WithRun(id.FromString("icc"), compilation.RunResult{Result: compilation.Result{Status: status.Flagged}}),
	)
}

// MockTimeoutRun expands to a realistic looking subject that contains some timed-out runs.
func MockTimeoutRun(name string) *subject.Subject {
	return subject.NewOrPanic(
		litmus.New("baz.litmus", litmus.WithThreads(4)),
		subject.WithRecipe(id.ArchPPC, MockRecipe("ppc")),
		subject.WithCompile(id.FromString("msvc"), MockSuccessfulCompile("msvc", name)),
		subject.WithRun(id.FromString("msvc"), compilation.RunResult{Result: compilation.Result{Status: status.RunTimeout}}),
	)
}

// MockSuccessfulCompile generates a mock CompileResult for a successful compile of subject sname with compiler cstr.
func MockSuccessfulCompile(cstr string, sname string) compilation.CompileResult {
	return compilation.CompileResult{
		Result: compilation.Result{
			Duration: 200 * time.Second,
			Status:   status.Ok,
		},
		Files: compilation.CompileFileset{
			Bin: path.Join(cstr, sname, "a.out"),
			Log: path.Join(cstr, sname, "log.txt"),
		},
	}
}

// MockRecipe constructs a mock recipe at dir.
func MockRecipe(dir string) recipe.Recipe {
	return recipe.New(
		dir,
		recipe.AddFiles("run.c", "aux.c", "aux.h"),
		recipe.CompileAllCToExe(),
	)
}
