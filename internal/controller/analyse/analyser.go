// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package analysis handles analysing a Corpus and filing its subjects into categorised sub-corpora.
package analyse

import (
	"context"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Analyser oversees the analysis of a particular plan.
type Analyser struct {
	// analysis is the analysis being built.
	analysis *analysis.Analysis

	// compilerTimes contains raw durations from each compiler's compilations.
	compilerTimes map[string][]time.Duration

	// runTimes contains raw durations from each compiler's runs.
	runTimes map[string][]time.Duration

	// corpus is the incoming corpus.
	corpus corpus.Corpus

	// nworkers is the number of workers.
	nworkers int
}

func (a *Analyser) Analyse(ctx context.Context) (*analysis.Analysis, error) {
	l := len(a.corpus)
	if l == 0 {
		return a.analysis, nil
	}

	if err := a.analyseCorpus(ctx, l); err != nil {
		return nil, err
	}
	for n, c := range a.analysis.Compilers {
		c.Time = analysis.NewTimeSet(a.compilerTimes[n]...)
		c.RunTime = analysis.NewTimeSet(a.runTimes[n]...)
		a.analysis.Compilers[n] = c
	}
	return a.analysis, nil
}

func (a *Analyser) analyseCorpus(ctx context.Context, l int) error {
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

// NewAnalyser initialises an analyser for plan p, with workers nworkers.
func NewAnalyser(p *plan.Plan, nworkers int) (*Analyser, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}

	lc := len(p.Compilers)
	a := Analyser{
		analysis: &analysis.Analysis{
			Plan:      p,
			ByStatus:  make(map[status.Status]corpus.Corpus, status.Last),
			Compilers: make(map[string]analysis.Compiler, lc),
		},
		corpus:        p.Corpus,
		compilerTimes: make(map[string][]time.Duration, lc),
		runTimes:      make(map[string][]time.Duration, lc),
		nworkers:      nworkers,
	}
	a.initCorpora(len(p.Corpus))
	a.initCompilers(p.Compilers)
	return &a, nil
}

func (a *Analyser) initCorpora(size int) {
	for i := status.Ok; i <= status.Last; i++ {
		a.analysis.ByStatus[i] = make(corpus.Corpus, size)
	}
}

func (a *Analyser) initCompilers(cs map[string]compiler.Compiler) {
	for cn, c := range cs {
		a.analysis.Compilers[cn] = analysis.Compiler{Counts: map[status.Status]int{}, Info: c}
		a.compilerTimes[cn] = []time.Duration{}
		a.runTimes[cn] = []time.Duration{}
	}
}

func classifyAndSend(named subject.Named, ch chan<- classification) {
	ch <- classify(named)
}

func (a *Analyser) build(ctx context.Context, ch <-chan classification, count int) error {
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

func (a *Analyser) Apply(r classification) {
	a.analysis.Flags |= r.flags
	for i := status.Ok; i <= status.Last; i++ {
		sf := i.Flag()

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
