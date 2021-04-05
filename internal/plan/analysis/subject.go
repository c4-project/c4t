// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"fmt"
	"time"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/timing"

	"github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/subject"
)

// subjectAnalysis holds the result of performing a single analysis on one subject.
type subjectAnalysis struct {
	flags        status.Flag
	sub          subject.Named
	cflags       map[id.ID]status.Flag
	ctimes       map[id.ID][]time.Duration
	clogs        map[id.ID]string
	rtimes       map[id.ID][]time.Duration
	cspan, rspan timing.Span
}

func newSubjectAnalysis(s subject.Named) subjectAnalysis {
	return subjectAnalysis{
		flags:  0,
		cflags: map[id.ID]status.Flag{},
		clogs:  map[id.ID]string{},
		ctimes: map[id.ID][]time.Duration{},
		rtimes: map[id.ID][]time.Duration{},
		sub:    s,
	}
}

// analyseSubject analyses the named subject s, using the compiler information ccs.
func (a *analyser) analyseSubject(s subject.Named) subjectAnalysis {
	c := newSubjectAnalysis(s)
	c.classifyCompilations(s.Compilations, a.analysis.Plan.Compilers, a.filters)
	return c
}

func (c *subjectAnalysis) classifyCompilations(crs compilation.Map, ccs compiler.InstanceMap, fs FilterSet) {
	for cid, cm := range crs {
		conf := ccs[cid]

		if cm.Compile != nil {
			c.classifyCompiler(cid, cm.Compile, conf, fs)
		}
		if cm.Run != nil {
			c.classifyRun(cid, cm.Run)
		}
	}
}

func (c *subjectAnalysis) classifyCompiler(cid id.ID, cm *compilation.CompileResult, conf compiler.Instance, fs FilterSet) {
	c.clogs[cid] = c.compileLog(cm)
	st, err := fs.FilteredStatus(cm.Status, conf, c.clogs[cid])
	if err != nil {
		// TODO(@MattWindsor91): do something about this!!
		return
	}
	c.logCompileStatus(cid, st)

	c.cspan.Union(cm.Timespan)
	if d := cm.Timespan.Duration(); d != 0 && st.CountsForTiming() {
		c.ctimes[cid] = append(c.ctimes[cid], d)
	}
}

func (c *subjectAnalysis) compileLog(cm *compilation.CompileResult) string {
	log, err := cm.Files.ReadLog("")
	if err != nil {
		return fmt.Sprintf("(ERROR GETTING COMPILE LOG: %s)", err)
	}
	return string(log)
}

func (c *subjectAnalysis) classifyRun(cid id.ID, r *compilation.RunResult) {
	// If we've already filtered this run out, don't unfilter it.
	if !(c.cflags[cid].MatchesStatus(status.Filtered)) {
		c.logCompileStatus(cid, r.Status)
	}

	c.rspan.Union(r.Timespan)
	if d := r.Timespan.Duration(); d != 0 && r.Status.CountsForTiming() {
		c.rtimes[cid] = append(c.rtimes[cid], d)
	}
}

func (c *subjectAnalysis) logCompileStatus(cid id.ID, s status.Status) {
	sf := s.Flag()
	c.flags |= sf
	c.cflags[cid] |= sf
}
