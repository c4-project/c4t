// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"fmt"
	"time"

	"github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/subject"
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
	c.classifyCompilations(s.Compilations, a.analysis.Plan.Compilers, a.filters)
	return c
}

func (c *subjectAnalysis) classifyCompilations(crs map[string]compilation.Compilation, ccs map[string]compiler.Instance, fs FilterSet) {
	for n, cm := range crs {
		conf := ccs[n]

		if cm.Compile != nil {
			c.classifyCompiler(n, cm.Compile, conf, fs)
		}
		if cm.Run != nil {
			c.classifyRun(n, cm.Run)
		}
	}
}

func (c *subjectAnalysis) classifyCompiler(cidstr string, cm *compilation.CompileResult, conf compiler.Instance, fs FilterSet) {
	c.clogs[cidstr] = c.compileLog(cm)
	st, err := fs.FilteredStatus(cm.Status, conf, c.clogs[cidstr])
	if err != nil {
		// TODO(@MattWindsor91): do something about this!!
		return
	}
	c.logCompileStatus(st, cidstr)

	if cm.Duration != 0 && st.CountsForTiming() {
		c.ctimes[cidstr] = append(c.ctimes[cidstr], cm.Duration)
	}
}

func (c *subjectAnalysis) compileLog(cm *compilation.CompileResult) string {
	log, err := cm.Files.ReadLog("")
	if err != nil {
		return fmt.Sprintf("(ERROR GETTING COMPILE LOG: %s)", err)
	}
	return string(log)
}

func (c *subjectAnalysis) classifyRun(cidstr string, r *compilation.RunResult) {
	// If we've already filtered this run out, don't unfilter it.
	if !(c.cflags[cidstr].MatchesStatus(status.Filtered)) {
		c.logCompileStatus(r.Status, cidstr)
	}
	if r.Duration != 0 && r.Status.CountsForTiming() {
		c.rtimes[cidstr] = append(c.rtimes[cidstr], r.Duration)
	}
}

func (c *subjectAnalysis) logCompileStatus(s status.Status, cidstr string) {
	sf := s.Flag()
	c.flags |= sf
	c.cflags[cidstr] |= sf
}
