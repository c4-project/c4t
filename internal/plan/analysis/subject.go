// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"time"

	"github.com/MattWindsor91/act-tester/internal/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/subject/status"

	"github.com/MattWindsor91/act-tester/internal/subject"
)

// subjectAnalysis holds the result of performing a single analysis on one subject.
type subjectAnalysis struct {
	flags  status.Flag
	sub    subject.Named
	cflags map[string]status.Flag
	ctimes map[string][]time.Duration
	rtimes map[string][]time.Duration
}

// analyseSubject analyses the named subject s.
func analyseSubject(s subject.Named) subjectAnalysis {
	c := subjectAnalysis{
		flags:  0,
		cflags: map[string]status.Flag{},
		ctimes: map[string][]time.Duration{},
		rtimes: map[string][]time.Duration{},
		sub:    s,
	}
	c.classifyCompiles(s.Compiles)
	c.classifyRuns(s.Runs)
	return c
}

func (c *subjectAnalysis) classifyCompiles(cs map[string]compilation.CompileResult) {
	for n, cm := range cs {
		sf := cm.Status.Flag()
		c.flags |= sf
		c.cflags[n] |= sf

		if cm.Duration != 0 && !(status.FlagFail | status.FlagTimeout).MatchesStatus(cm.Status) {
			c.ctimes[n] = append(c.ctimes[n], cm.Duration)
		}
	}
}

func (c *subjectAnalysis) classifyRuns(rs map[string]compilation.RunResult) {
	for n, r := range rs {
		sf := r.Status.Flag()
		c.flags |= sf
		c.cflags[n] |= sf

		if r.Duration != 0 && !(status.FlagRunFail | status.FlagRunTimeout).MatchesStatus(r.Status) {
			c.rtimes[n] = append(c.rtimes[n], r.Duration)
		}
	}
}
