// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package analysis handles analysing a Corpus and filing its subjects into categorised sub-corpi.
package analysis

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/model/compiler"
	"github.com/MattWindsor91/act-tester/internal/model/plan"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Analyse analyses a plan p using up to nworkers workers.
func Analyse(ctx context.Context, p *plan.Plan, nworkers int) (*Analysis, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}

	l := len(p.Corpus)

	col := initAnalysis(l, p.Compilers)
	// Various bits of the collator expect there to be at least one subject.
	if l == 0 {
		return col, nil
	}

	ch := make(chan classification)
	err := p.Corpus.Par(ctx, nworkers,
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

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	return p.Check()
}

func initAnalysis(l int, compilers map[string]compiler.Compiler) *Analysis {
	col := Analysis{
		ByStatus:  make(map[subject.Status]corpus.Corpus, subject.NumStatus),
		Compilers: make(map[string]Compiler, len(compilers)),
	}
	for i := subject.StatusOk; i < subject.NumStatus; i++ {
		col.ByStatus[i] = make(corpus.Corpus, l)
	}
	for cn, c := range compilers {
		col.Compilers[cn] = Compiler{Counts: map[subject.Status]int{}, Info: c}
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
	a.Flags |= r.flags
	for i := subject.StatusOk; i < subject.NumStatus; i++ {
		sf := statusFlags[i]

		if r.flags.matches(sf) {
			a.ByStatus[i][r.sub.Name] = r.sub.Subject
		}

		for cstr, f := range r.compilers {
			if _, ok := a.Compilers[cstr]; !ok {
				continue
			}

			if f.matches(sf) {
				a.Compilers[cstr].Counts[i]++
			}
		}
	}

}

type classification struct {
	flags     Flag
	compilers map[string]Flag
	sub       subject.Named
}
