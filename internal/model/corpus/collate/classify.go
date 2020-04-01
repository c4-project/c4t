// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate

import "github.com/MattWindsor91/act-tester/internal/model/subject"

func classify(named subject.Named) flag {
	var f flag
	f |= classifyCompiles(named.Compiles)
	f |= classifyRuns(named.Runs)
	return f
}

func classifyCompiles(cs map[string]subject.CompileResult) flag {
	f := flagOk
	for _, c := range cs {
		f |= classifyStatus(c.Status)
	}
	return f
}

func classifyRuns(rs map[string]subject.RunResult) flag {
	f := flagOk
	for _, r := range rs {
		f |= classifyStatus(r.Status)
	}
	return f
}

func classifyStatus(s subject.Status) flag {
	switch s {
	case subject.StatusFlagged:
		return flagFlagged
	case subject.StatusCompileTimeout:
		return flagCompileTimeout
	case subject.StatusCompileFail:
		return flagCompileFail
	case subject.StatusRunTimeout:
		return flagRunTimeout
	case subject.StatusRunFail:
		return flagRunFailure
	default:
		return flagOk
	}
}
