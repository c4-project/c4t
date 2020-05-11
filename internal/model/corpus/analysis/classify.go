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

func classify(named subject.Named) classification {
	c := classification{
		flags:     flagOk,
		compilers: map[string]flag{},
		sub:       named,
	}
	c.classifyCompiles(named.Compiles)
	c.classifyRuns(named.Runs)
	return c
}

func (c *classification) classifyCompiles(cs map[string]subject.CompileResult) {
	for n, cm := range cs {
		c.compilers[n] |= statusFlags[cm.Status]
		c.flags |= statusFlags[cm.Status]
	}
}

func (c *classification) classifyRuns(rs map[string]subject.RunResult) {
	for n, r := range rs {
		c.compilers[n] |= statusFlags[r.Status]
		c.flags |= statusFlags[r.Status]
	}
}
