// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package analysis handles analysing a Corpus and filing its subjects into categorised sub-corpi.
package analysis

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Analyse collates a corpus c using up to nworkers workers.
func Analyse(ctx context.Context, c corpus.Corpus, nworkers int) (*Analysis, error) {
	l := len(c)
	col := Analysis{
		ByStatus: make(map[subject.Status]corpus.Corpus, subject.NumStatus),
	}
	for i := subject.StatusOk; i < subject.NumStatus; i++ {
		col.ByStatus[i] = make(corpus.Corpus, l)
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

func (c *Analysis) build(ctx context.Context, ch <-chan collationRequest, count int) error {
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

func (c *Analysis) file(rq collationRequest) {
	for s, c := range c.ByStatus {
		if rq.flags.matches(statusFlags[s]) {
			c[rq.sub.Name] = rq.sub.Subject
		}
	}
}

type collationRequest struct {
	flags flag
	sub   subject.Named
}
