// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus

import (
	"path"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Mock produces a representative corpus including the following features:
// - a subject with a failed compilation;
// - a subject with a flagged observation.
func Mock() Corpus {
	return Corpus{
		"foo":    subject.Subject{Stats: model.Statset{Threads: 1}, OrigLitmus: "foo.litmus"},
		"bar":    MockFailedCompile("bar"),
		"baz":    MockFlaggedRun("baz"),
		"barbaz": MockTimeoutRun("barbaz"),
	}
}

// MockFailedCompile expands to a realistic looking Subject that contains a failed compilation.
func MockFailedCompile(name string) subject.Subject {
	return subject.Subject{
		Stats: model.Statset{
			Threads: 8,
		},
		OrigLitmus: name + ".litmus",
		Harnesses: map[string]subject.Harness{
			id.ArchArm.String(): {
				Dir:   "arm",
				Files: []string{"run.c", "aux.c", "aux.h"},
			},
		},
		Compiles: map[string]subject.CompileResult{
			"gcc": {
				Result: subject.Result{Status: status.CompileFail},
				Files:  subject.CompileFileset{},
			},
			"clang": MockSuccessfulCompile("clang", name),
		},
		Runs: map[string]subject.RunResult{
			"gcc": {
				Result: subject.Result{Status: status.CompileFail},
			},
			"clang": {
				Result: subject.Result{Status: status.Ok},
			},
		},
	}
}

// MockFlaggedRun expands to a realistic looking Subject that contains some flagged runs.
func MockFlaggedRun(name string) subject.Subject {
	return subject.Subject{
		Stats:      model.Statset{Threads: 2},
		OrigLitmus: name + ".litmus",
		Harnesses: map[string]subject.Harness{
			id.ArchX8664.String(): MockHarness("x86"),
		},
		Compiles: map[string]subject.CompileResult{
			"gcc": MockSuccessfulCompile("gcc", name),
			"icc": MockSuccessfulCompile("icc", name),
		},
		Runs: map[string]subject.RunResult{
			"gcc": {Result: subject.Result{Status: status.Flagged}},
			"icc": {Result: subject.Result{Status: status.Flagged}},
		},
	}
}

// MockTimeoutRun expands to a realistic looking Subject that contains some timed-out runs.
func MockTimeoutRun(name string) subject.Subject {
	return subject.Subject{
		Stats:      model.Statset{Threads: 4},
		OrigLitmus: "baz.litmus",
		Harnesses: map[string]subject.Harness{
			id.ArchX8664.String(): MockHarness("x86"),
			id.ArchPPC.String():   MockHarness("ppc"),
		},
		Compiles: map[string]subject.CompileResult{
			"msvc": MockSuccessfulCompile("msvc", name),
		},
		Runs: map[string]subject.RunResult{
			"msvc": {Result: subject.Result{Status: status.RunTimeout}},
		},
	}
}

// MockSuccessfulCompile generates a mock CompileResult for a successful compile of subject sname with compiler cstr.
func MockSuccessfulCompile(cstr string, sname string) subject.CompileResult {
	return subject.CompileResult{
		Result: subject.Result{
			Duration: 200 * time.Second,
			Status:   status.Ok,
		},
		Files: subject.CompileFileset{
			Bin: path.Join(cstr, sname, "a.out"),
			Log: path.Join(cstr, sname, "log.txt"),
		},
	}
}

// MockHarness constructs a mock harness at dir.
func MockHarness(dir string) subject.Harness {
	return subject.Harness{
		Dir:   dir,
		Files: []string{"run.c", "aux.c", "aux.h"},
	}
}
