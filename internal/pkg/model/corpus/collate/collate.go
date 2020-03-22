// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// collate handles analysing a Corpus and filing its subjects into categorised sub-corpi
package collate

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"
)

// Collation represents a grouping of corpus subjects according to various issues.
type Collation struct {
	// CompileFailures contains the subset of the collated corpus that ran into compiler failures.
	CompileFailures corpus.Corpus
	// Flagged contains the subset of the collated corpus that has been flagged as possibly buggy.
	Flagged corpus.Corpus
	// RunFailures contains the subset of the collated corpus that ran into runtime failures.
	RunFailures corpus.Corpus
	// Timeouts contains the subset of the collated corpus that timed out.
	Timeouts corpus.Corpus
	// Successes contains the subset of the collated corpus that doesn't fit into any of the above boxes.
	Successes corpus.Corpus
}

// Collate collates a corpus c using up to nworkers workers.
func Collate(ctx context.Context, c corpus.Corpus, nworkers int) (*Collation, error) {
	l := len(c)
	col := Collation{
		CompileFailures: make(corpus.Corpus, l),
		RunFailures:     make(corpus.Corpus, l),
		Timeouts:        make(corpus.Corpus, l),
		Flagged:         make(corpus.Corpus, l),
		Successes:       make(corpus.Corpus, l),
	}

	// Various bits of the collator expect there to be at least one subject.
	if l == 0 {
		return &col, nil
	}

	ch := make(chan collationRequest)
	err := c.Par(ctx, nworkers,
		func(ctx context.Context, named subject.Named) error {
			classifyAndSend(named, ch)
			return nil
		},
		func(ctx context.Context) error {
			return col.build(ctx, ch, l)
		},
	)
	return &col, err
}

func classifyAndSend(named subject.Named, ch chan<- collationRequest) {
	fs := classify(named)
	ch <- collationRequest{
		flags: fs,
		sub:   named,
	}
}

func (c *Collation) build(ctx context.Context, ch <-chan collationRequest, count int) error {
	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case rq := <-ch:
			c.file(rq)
		}
	}
	return nil
}

func (c *Collation) file(rq collationRequest) {
	cases := map[collationFlag]corpus.Corpus{
		ccOk:         c.Successes,
		ccCompile:    c.CompileFailures,
		ccRunFailure: c.RunFailures,
		ccFlag:       c.Flagged,
		ccTimeout:    c.Timeouts,
	}
	for f, c := range cases {
		if rq.flags.matches(f) {
			c[rq.sub.Name] = rq.sub.Subject
		}
	}
}

type collationRequest struct {
	flags collationFlag
	sub   subject.Named
}
