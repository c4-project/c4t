// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

type classification struct {
	flags  Flag
	sub    subject.Named
	cflags map[string]Flag
	ctimes map[string][]time.Duration
	rtimes map[string][]time.Duration
}

func classify(named subject.Named) classification {
	c := classification{
		flags:  FlagOk,
		cflags: map[string]Flag{},
		ctimes: map[string][]time.Duration{},
		rtimes: map[string][]time.Duration{},
		sub:    named,
	}
	c.classifyCompiles(named.Compiles)
	c.classifyRuns(named.Runs)
	return c
}

func (c *classification) classifyCompiles(cs map[string]subject.CompileResult) {
	for n, cm := range cs {
		sf := statusFlags[cm.Status]
		c.flags |= sf
		c.cflags[n] |= sf

		if cm.Duration != 0 && !(FlagFail | FlagTimeout).Matches(sf) {
			c.ctimes[n] = append(c.ctimes[n], cm.Duration)
		}
	}
}

func (c *classification) classifyRuns(rs map[string]subject.RunResult) {
	for n, r := range rs {
		sf := statusFlags[r.Status]
		c.flags |= sf
		c.cflags[n] |= sf

		if r.Duration != 0 && !(FlagRunFail | FlagRunTimeout).Matches(sf) {
			c.rtimes[n] = append(c.rtimes[n], r.Duration)
		}
	}
}
