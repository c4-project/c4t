// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

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

// analyseSubject analyses the named subject s, using the compiler information ccs.
func (a *analyser) analyseSubject(s subject.Named) subjectAnalysis {
	c := newSubjectAnalysis(s)
	c.classifyCompiles(s.Compiles, a.analysis.Plan.Compilers, a.filters)
	c.classifyRuns(s.Runs)
	return c
}

func (c *subjectAnalysis) classifyCompiles(crs map[string]compilation.CompileResult, ccs map[string]compiler.Configuration, fs FilterSet) {
	for n, cm := range crs {
		c.clogs[n] = c.compileLog(cm)
		st, err := fs.FilteredStatus(cm, ccs[n], c.clogs[n])
		if err != nil {
			// TODO(@MattWindsor91): do something about this!!
			return
		}
		c.logCompileStatus(st, n)

		if cm.Duration != 0 && !(status.FlagFiltered | status.FlagFail | status.FlagTimeout).MatchesStatus(st) {
			c.ctimes[n] = append(c.ctimes[n], cm.Duration)
		}
	}
}

func (c *subjectAnalysis) compileLog(cm compilation.CompileResult) string {
	log, err := cm.Files.ReadLog("")
	if err != nil {
		return fmt.Sprintf("(ERROR GETTING COMPILE LOG: %s)", err)
	}
	return string(log)
}

func (c *subjectAnalysis) classifyRuns(rs map[string]compilation.RunResult) {
	for n, r := range rs {
		c.logCompileStatus(r.Status, n)

		if r.Duration != 0 && !(status.FlagRunFail | status.FlagRunTimeout).MatchesStatus(r.Status) {
			c.rtimes[n] = append(c.rtimes[n], r.Duration)
		}
	}
}

func (c *subjectAnalysis) logCompileStatus(s status.Status, cidstr string) {
	sf := s.Flag()
	c.flags |= sf
	c.cflags[cidstr] |= sf
}
