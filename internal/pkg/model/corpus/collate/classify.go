// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate

import "github.com/MattWindsor91/act-tester/internal/pkg/model/subject"

func classify(named subject.Named) flag {
	var f flag
	f |= classifyCompiles(named.Compiles)
	f |= classifyRuns(named.Runs)
	return f
}

func classifyCompiles(cs map[string]subject.CompileResult) flag {
	for _, c := range cs {
		if !c.Success {
			return flagCompileFail
		}
	}
	return flagOk
}

func classifyRuns(rs map[string]subject.Run) flag {
	f := flagOk
	for _, r := range rs {
		f |= classifyRunStatus(r.Status)
	}
	return f
}

func classifyRunStatus(s subject.Status) flag {
	switch s {
	case subject.StatusFlagged:
		return flagFlagged
	case subject.StatusCompileTimeout:
		// This (will be) probably redundant, as we classify the compile too.
		return flagCompileTimeout
	case subject.StatusCompileFail:
		// This is probably redundant, as we classify the compile too.
		return flagCompileFail
	case subject.StatusRunTimeout:
		return flagRunTimeout
	case subject.StatusRunFail:
		return flagRunFailure
	default:
		return flagOk
	}
}
