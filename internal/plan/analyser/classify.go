// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analyser

import (
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

type classification struct {
	flags  status.Flag
	sub    subject.Named
	cflags map[string]status.Flag
	ctimes map[string][]time.Duration
	rtimes map[string][]time.Duration
}

func classify(named subject.Named) classification {
	c := classification{
		flags:  status.FlagOk,
		cflags: map[string]status.Flag{},
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
		sf := cm.Status.Flag()
		c.flags |= sf
		c.cflags[n] |= sf

		if cm.Duration != 0 && !(status.FlagFail | status.FlagTimeout).Matches(sf) {
			c.ctimes[n] = append(c.ctimes[n], cm.Duration)
		}
	}
}

func (c *classification) classifyRuns(rs map[string]subject.RunResult) {
	for n, r := range rs {
		sf := r.Status.Flag()
		c.flags |= sf
		c.cflags[n] |= sf

		if r.Duration != 0 && !(status.FlagRunFail | status.FlagRunTimeout).Matches(sf) {
			c.rtimes[n] = append(c.rtimes[n], r.Duration)
		}
	}
}
