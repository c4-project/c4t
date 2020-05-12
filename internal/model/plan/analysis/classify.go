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
	cflags map[string]Flag
	ctimes map[string][]time.Duration
	sub    subject.Named
}

func classify(named subject.Named) classification {
	c := classification{
		flags:  FlagOk,
		cflags: map[string]Flag{},
		ctimes: map[string][]time.Duration{},
		sub:    named,
	}
	c.classifyCompiles(named.Compiles)
	c.classifyRuns(named.Runs)
	return c
}

func (c *classification) classifyCompiles(cs map[string]subject.CompileResult) {
	for n, cm := range cs {
		c.cflags[n] |= statusFlags[cm.Status]
		c.flags |= statusFlags[cm.Status]
		if cm.Duration != 0 {
			c.ctimes[n] = append(c.ctimes[n], cm.Duration)
		}
	}
}

func (c *classification) classifyRuns(rs map[string]subject.RunResult) {
	for n, r := range rs {
		c.cflags[n] |= statusFlags[r.Status]
		c.flags |= statusFlags[r.Status]
	}
}
