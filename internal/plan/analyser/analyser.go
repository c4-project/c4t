// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package analyser handles analysing a Corpus and filing its subjects into categorised sub-corpora.
package analyser

import (
	"context"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/MattWindsor91/act-tester/internal/plan"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Analyser oversees the analyser of a particular plan.
type Analyser struct {
	// analyser is the analyser being built.
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

// Analyse runs the analyser with context ctx.
func (a *Analyser) Analyse(ctx context.Context) (*Analysis, error) {
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

// New initialises an analyser for plan p, with workers nworkers.
func New(p *plan.Plan, nworkers int) (*Analyser, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}

	lc := len(p.Compilers)
	a := Analyser{
		analysis:      newAnalysis(p),
		corpus:        p.Corpus,
		compilerTimes: make(map[string][]time.Duration, lc),
		runTimes:      make(map[string][]time.Duration, lc),
		nworkers:      nworkers,
	}
	a.initCompilers(p.Compilers)
	return &a, nil
}

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	return p.Check()
}

func (a *Analyser) initCompilers(cs map[string]compiler.Compiler) {
	for cn, c := range cs {
		a.analysis.Compilers[cn] = Compiler{Counts: map[status.Status]int{}, Info: c}
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
			a.apply(rq)
		}
	}
	return nil
}

func (a *Analyser) apply(r classification) {
	a.analysis.Flags |= r.flags
	for i := status.Ok; i <= status.Last; i++ {
		sf := i.Flag()
		a.applyByStatus(i, sf, r)
		a.applyCompilers(i, sf, r)
		for cstr, ts := range r.ctimes {
			a.compilerTimes[cstr] = append(a.compilerTimes[cstr], ts...)
		}
		for cstr, ts := range r.rtimes {
			a.runTimes[cstr] = append(a.runTimes[cstr], ts...)
		}
	}
}

func (a *Analyser) applyByStatus(s status.Status, sf status.Flag, r classification) {
	if !r.flags.Matches(sf) {
		return
	}
	if _, ok := a.analysis.ByStatus[s]; !ok {
		a.analysis.ByStatus[s] = make(corpus.Corpus)
	}
	a.analysis.ByStatus[s][r.sub.Name] = r.sub.Subject
}

func (a *Analyser) applyCompilers(s status.Status, sf status.Flag, r classification) {
	for cstr, f := range r.cflags {
		a.applyCompiler(s, sf, f, cstr)
	}
}

func (a *Analyser) applyCompiler(s status.Status, sf, cf status.Flag, cstr string) {
	if !cf.Matches(sf) {
		return
	}
	if _, ok := a.analysis.Compilers[cstr]; !ok {
		return
	}
	a.analysis.Compilers[cstr].Counts[s]++
}
