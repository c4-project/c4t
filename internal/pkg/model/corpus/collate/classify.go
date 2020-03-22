// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate

import "github.com/MattWindsor91/act-tester/internal/pkg/model/subject"

func classify(named subject.Named) collationFlag {
	var f collationFlag
	f |= classifyCompiles(named.Compiles)
	f |= classifyRuns(named.Runs)
	return f
}

func classifyCompiles(cs map[string]subject.CompileResult) collationFlag {
	for _, c := range cs {
		if !c.Success {
			return ccCompile
		}
	}
	return ccOk
}

func classifyRuns(rs map[string]subject.Run) collationFlag {
	f := ccOk
	for _, r := range rs {
		f |= classifyRun(r.Status)
	}
	return f
}

func classifyRun(s subject.Status) collationFlag {
	switch s {
	case subject.StatusFlagged:
		return ccFlag
	case subject.StatusCompileFail:
		// This is probably redundant, as we classify the compile too.
		return ccCompile
	case subject.StatusTimeout:
		return ccTimeout
	case subject.StatusUnknown:
		// TODO(@MattWindsor91): run failure
		return ccRunFailure
	default:
		return ccOk
	}
}
