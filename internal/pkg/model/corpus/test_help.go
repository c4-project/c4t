// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus

import (
	"path"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"
	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"
)

// Mock produces a representative corpus including the following features:
// - a subject with a failed compilation;
// - a subject with a flagged observation.
func Mock() Corpus {
	return Corpus{
		"foo":    subject.Subject{Threads: 1, Litmus: "foo.litmus"},
		"bar":    MockFailedCompile("bar"),
		"baz":    MockFlaggedRun("baz"),
		"barbaz": MockTimeoutRun("barbaz"),
	}
}

// MockFailedCompile expands to a realistic looking Subject that contains a failed compilation.
func MockFailedCompile(name string) subject.Subject {
	return subject.Subject{
		Threads: 8,
		Litmus:  name + ".litmus",
		Harnesses: map[string]subject.Harness{
			id.ArchArm.String(): {
				Dir:   "arm",
				Files: []string{"run.c", "aux.c", "aux.h"},
			},
		},
		Compiles: map[string]subject.CompileResult{
			"gcc": {
				Success: false,
				Files:   subject.CompileFileset{},
			},
			"clang": MockSuccessfulCompile("clang", name),
		},
		Runs: map[string]subject.Run{
			"gcc": {
				Status: subject.StatusOk,
			},
			"clang": {
				Status: subject.StatusCompileFail,
			},
		},
	}
}

// MockFlaggedRun expands to a realistic looking Subject that contains some flagged runs.
func MockFlaggedRun(name string) subject.Subject {
	return subject.Subject{
		Threads: 2,
		Litmus:  name + ".litmus",
		Harnesses: map[string]subject.Harness{
			id.ArchX8664.String(): MockHarness("x86"),
		},
		Compiles: map[string]subject.CompileResult{
			"gcc": MockSuccessfulCompile("gcc", name),
			"icc": MockSuccessfulCompile("icc", name),
		},
		Runs: map[string]subject.Run{
			"gcc": {Status: subject.StatusFlagged},
			"icc": {Status: subject.StatusFlagged},
		},
	}
}

// MockTimeoutRun expands to a realistic looking Subject that contains some timed-out runs.
func MockTimeoutRun(name string) subject.Subject {
	return subject.Subject{
		Threads: 4,
		Litmus:  "baz.litmus",
		Harnesses: map[string]subject.Harness{
			id.ArchX8664.String(): MockHarness("x86"),
			id.ArchPPC.String():   MockHarness("ppc"),
		},
		Compiles: map[string]subject.CompileResult{
			"msvc": MockSuccessfulCompile("msvc", name),
		},
		Runs: map[string]subject.Run{
			"msvc": {Status: subject.StatusTimeout},
		},
	}
}

// MockSuccessfulCompile generates a mock CompileResult for a successful compile of subject sname with compiler cstr.
func MockSuccessfulCompile(cstr string, sname string) subject.CompileResult {
	return subject.CompileResult{
		Success: true,
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
