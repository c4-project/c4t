// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/id"

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
	clogs  map[string]string
	rtimes map[string][]time.Duration
}

func newSubjectAnalysis(s subject.Named) subjectAnalysis {
	return subjectAnalysis{
		flags:  0,
		cflags: map[string]status.Flag{},
		clogs:  map[string]string{},
		ctimes: map[string][]time.Duration{},
		rtimes: map[string][]time.Duration{},
		sub:    s,
	}
}

// analyseSubject analyses the named subject s.
func analyseSubject(s subject.Named) subjectAnalysis {
	c := newSubjectAnalysis(s)
	c.classifyCompiles(s.Compiles)
	c.classifyRuns(s.Runs)
	return c
}

func (c *subjectAnalysis) classifyCompiles(cs map[string]compilation.CompileResult) {
	for n, cm := range cs {
		c.logCompileFlag(cm.Status.Flag(), n)
		c.clogs[n] = c.compileLog(n)

		if cm.Duration != 0 && !(status.FlagFail | status.FlagTimeout).MatchesStatus(cm.Status) {
			c.ctimes[n] = append(c.ctimes[n], cm.Duration)
		}
	}
}

func (c *subjectAnalysis) compileLog(cidstr string) string {
	log, err := c.tryCompileLog(cidstr)
	if err != nil {
		return fmt.Sprintf("(ERROR GETTING COMPILE LOG: %s)", err)
	}
	return log
}

func (c *subjectAnalysis) tryCompileLog(cidstr string) (string, error) {
	cid, err := id.TryFromString(cidstr)
	if err != nil {
		return "", err
	}
	bs, err := c.sub.ReadCompilerLog("", cid)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (c *subjectAnalysis) classifyRuns(rs map[string]compilation.RunResult) {
	for n, r := range rs {
		c.logCompileFlag(r.Status.Flag(), n)

		if r.Duration != 0 && !(status.FlagRunFail | status.FlagRunTimeout).MatchesStatus(r.Status) {
			c.rtimes[n] = append(c.rtimes[n], r.Duration)
		}
	}
}

func (c *subjectAnalysis) logCompileFlag(sf status.Flag, cidstr string) {
	c.flags |= sf
	c.cflags[cidstr] |= sf
}
