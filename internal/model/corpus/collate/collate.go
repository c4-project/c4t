// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package collate handles analysing a Corpus and filing its subjects into categorised sub-corpi.
package collate

import (
	"context"
	"fmt"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Collation represents a grouping of corpus subjects according to various issues.
type Collation struct {
	// ByStatus maps each status to the corpus of subjects that fall into it.
	ByStatus map[subject.Status]corpus.Corpus
}

// String summarises this collation as a string.
func (c *Collation) String() string {
	var sb strings.Builder

	bf := c.ByStatus

	// We range over this to enforce a deterministic order.
	for i := subject.StatusOk; i < subject.NumStatus; i++ {
		if i != subject.StatusOk {
			sb.WriteString(", ")
		}
		_, _ = fmt.Fprintf(&sb, "%d %s", len(bf[i]), i.String())
	}

	return sb.String()
}

// Collate collates a corpus c using up to nworkers workers.
func Collate(ctx context.Context, c corpus.Corpus, nworkers int) (*Collation, error) {
	l := len(c)
	col := Collation{
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
