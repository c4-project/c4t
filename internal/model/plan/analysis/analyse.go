// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package analysis handles analysing a Corpus and filing its subjects into categorised sub-corpi.
package analysis

import (
	"context"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/compiler"
	"github.com/MattWindsor91/act-tester/internal/model/plan"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

type analyser struct {
	// analysis is the analysis being built.
	analysis *Analysis

	// compilerTimes contains raw durations from each compiler's compilations.
	compilerTimes map[string][]time.Duration

	// runTimes contains raw durations from each compiler's runs.
	runTimes map[string][]time.Duration

	// corpus is the incoming corpus.
	corpus corpus.Corpus

	// nworkers is the number of workers.
	nworkers int
}

// Analyse analyses a plan p using up to nworkers workers.
func Analyse(ctx context.Context, p *plan.Plan, nworkers int) (*Analysis, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}

	return initAnalyser(p.Corpus, p.Compilers, nworkers).analyse(ctx)
}

func (a *analyser) analyse(ctx context.Context) (*Analysis, error) {
	l := len(a.corpus)
	if l == 0 {
		return a.analysis, nil
	}

	if err := a.analyseCorpus(ctx, l); err != nil {
		return nil, err
	}
	for n, c := range a.analysis.Compilers {
		c.Time = NewTimeSet(a.compilerTimes[n]...)
		c.RunTime = NewTimeSet(a.runTimes[n]...)
		a.analysis.Compilers[n] = c
	}
	return a.analysis, nil
}

func (a *analyser) analyseCorpus(ctx context.Context, l int) error {
	ch := make(chan classification)
	err := a.corpus.Par(ctx, a.nworkers,
		func(ctx context.Context, named subject.Named) error {
			classifyAndSend(named, ch)
			return nil
		},
		func(ctx context.Context) error {
			return a.build(ctx, ch, l)
		},
	)
	return err
}

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	return p.Check()
}

func initAnalyser(c corpus.Corpus, compilers map[string]compiler.Compiler, nworkers int) *analyser {
	lc := len(compilers)
	a := analyser{
		analysis: &Analysis{
			ByStatus:  make(map[status.Status]corpus.Corpus, status.Num),
			Compilers: make(map[string]Compiler, lc),
		},
		corpus:        c,
		compilerTimes: make(map[string][]time.Duration, lc),
		runTimes:      make(map[string][]time.Duration, lc),
		nworkers:      nworkers,
	}
	for i := status.Ok; i < status.Num; i++ {
		a.analysis.ByStatus[i] = make(corpus.Corpus, len(c))
	}
	for cn, c := range compilers {
		a.analysis.Compilers[cn] = Compiler{Counts: map[status.Status]int{}, Info: c}
		a.compilerTimes[cn] = []time.Duration{}
		a.runTimes[cn] = []time.Duration{}
	}
	return &a
}

func classifyAndSend(named subject.Named, ch chan<- classification) {
	ch <- classify(named)
}

func (a *analyser) build(ctx context.Context, ch <-chan classification, count int) error {
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

func (a *analyser) Apply(r classification) {
	a.analysis.Flags |= r.flags
	for i := status.Ok; i < status.Num; i++ {
		sf := statusFlags[i]

		if r.flags.Matches(sf) {
			a.analysis.ByStatus[i][r.sub.Name] = r.sub.Subject
		}

		for cstr, f := range r.cflags {
			if _, ok := a.analysis.Compilers[cstr]; !ok {
				continue
			}

			if f.Matches(sf) {
				a.analysis.Compilers[cstr].Counts[i]++
			}
		}
		for cstr, ts := range r.ctimes {
			a.compilerTimes[cstr] = append(a.compilerTimes[cstr], ts...)
		}
		for cstr, ts := range r.rtimes {
			a.runTimes[cstr] = append(a.runTimes[cstr], ts...)
		}
	}
}
