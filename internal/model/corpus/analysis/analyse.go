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

	col := initAnalysis(l)
	// Various bits of the collator expect there to be at least one subject.
	if l == 0 {
		return col, nil
	}

	ch := make(chan classification)
	err := c.Par(ctx, nworkers,
		func(ctx context.Context, named subject.Named) error {
			classifyAndSend(named, ch)
			return nil
		},
		func(ctx context.Context) error {
			return col.build(ctx, ch, l)
		},
	)
	return col, err
}

func initAnalysis(l int) *Analysis {
	col := Analysis{
		ByStatus:       make(map[subject.Status]corpus.Corpus, subject.NumStatus),
		CompilerCounts: map[string]map[subject.Status]int{},
	}
	for i := subject.StatusOk; i < subject.NumStatus; i++ {
		col.ByStatus[i] = make(corpus.Corpus, l)
	}
	return &col
}

func classifyAndSend(named subject.Named, ch chan<- classification) {
	ch <- classify(named)
}

func (a *Analysis) build(ctx context.Context, ch <-chan classification, count int) error {
	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case rq := <-ch:
			a.Apply(rq)
		}
	}
	return nil
}

func (a *Analysis) Apply(r classification) {
	for i := subject.StatusOk; i < subject.NumStatus; i++ {
		sf := statusFlags[i]

		if r.flags.matches(sf) {
			a.ByStatus[i][r.sub.Name] = r.sub.Subject
		}

		for cstr, f := range r.compilers {
			if f.matches(sf) {
				if a.CompilerCounts[cstr] == nil {
					a.CompilerCounts[cstr] = make(map[subject.Status]int)
				}
				a.CompilerCounts[cstr][i]++
			}
		}
	}

}

type classification struct {
	flags     flag
	compilers map[string]flag
	sub       subject.Named
}
