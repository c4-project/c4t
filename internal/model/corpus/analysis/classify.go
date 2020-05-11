// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import "github.com/MattWindsor91/act-tester/internal/model/subject"

var statusFlags = [subject.NumStatus]flag{
	subject.StatusOk:             flagOk,
	subject.StatusFlagged:        flagFlagged,
	subject.StatusCompileTimeout: flagCompileTimeout,
	subject.StatusCompileFail:    flagCompileFail,
	subject.StatusRunTimeout:     flagRunTimeout,
	subject.StatusRunFail:        flagRunFailure,
}

func classify(named subject.Named) flag {
	var f flag
	f |= classifyCompiles(named.Compiles)
	f |= classifyRuns(named.Runs)
	return f
}

func classifyCompiles(cs map[string]subject.CompileResult) flag {
	f := flagOk
	for _, c := range cs {
		f |= statusFlags[c.Status]
	}
	return f
}

func classifyRuns(rs map[string]subject.RunResult) flag {
	f := flagOk
	for _, r := range rs {
		f |= statusFlags[r.Status]
	}
	return f
}
