// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package analysis handles analysing a Corpus and filing its subjects into categorised sub-corpora.
package analysis

import (
	"context"
	"time"

	"github.com/MattWindsor91/act-tester/internal/subject/status"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/MattWindsor91/act-tester/internal/plan"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus"

	"github.com/MattWindsor91/act-tester/internal/subject"
)

// Analyse runs the analyser with context ctx, on plan p and with nworkers parallel subject analysers.
func Analyse(ctx context.Context, p *plan.Plan, nworkers int) (*Analysis, error) {
	a, err := New(p, nworkers)
	if err != nil {
		return nil, err
	}
	return a.analyse(ctx)
}

// analyser oversees the analysis of a particular plan.
//
// The analysis proceeds by classifying individual subjects with a degree of parallelism, and builds an
// Analysis in place by collating those classifications as they come in.
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

// analyse runs the analyser with context ctx.
func (a *analyser) analyse(ctx context.Context) (*Analysis, error) {
	if err := a.analyseCorpus(ctx); err != nil {
		return nil, err
	}
	for n, c := range a.analysis.Compilers {
		c.Time = NewTimeSet(a.compilerTimes[n]...)
		c.RunTime = NewTimeSet(a.runTimes[n]...)
		a.analysis.Compilers[n] = c
	}
	return a.analysis, nil
}

func (a *analyser) analyseCorpus(ctx context.Context) error {
	ch := make(chan classification)
	err := a.corpus.Par(ctx, a.nworkers,
		func(ctx context.Context, named subject.Named) error {
			classifyAndSend(named, ch)
			return nil
		},
		func(ctx context.Context) error {
			return a.build(ctx, ch)
		},
	)
	return err
}

// New initialises an analyser for plan p, with workers nworkers.
func New(p *plan.Plan, nworkers int) (*analyser, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}

	lc := len(p.Compilers)
	a := analyser{
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

func (a *analyser) initCompilers(cs map[string]compiler.Configuration) {
	for cn, c := range cs {
		a.analysis.Compilers[cn] = Compiler{Counts: map[status.Status]int{}, Info: c}
		a.compilerTimes[cn] = []time.Duration{}
		a.runTimes[cn] = []time.Duration{}
	}
}

func classifyAndSend(named subject.Named, ch chan<- classification) {
	ch <- classify(named)
}

func (a *analyser) build(ctx context.Context, ch <-chan classification) error {
	for i := 0; i < len(a.corpus); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case rq := <-ch:
			a.apply(rq)
		}
	}
	return nil
}

func (a *analyser) apply(r classification) {
	a.analysis.Flags |= r.flags
	for i := status.Ok; i <= status.Last; i++ {
		a.applyByStatus(i, r)
		a.applyCompilers(i, r)
		a.applyTimes(r)
	}
}

func (a *analyser) applyByStatus(s status.Status, r classification) {
	if !r.flags.MatchesStatus(s) {
		return
	}
	if _, ok := a.analysis.ByStatus[s]; !ok {
		a.analysis.ByStatus[s] = make(corpus.Corpus)
	}
	a.analysis.ByStatus[s][r.sub.Name] = r.sub.Subject
}

func (a *analyser) applyCompilers(s status.Status, r classification) {
	for cstr, f := range r.cflags {
		a.applyCompiler(s, f, cstr)
	}
}

func (a *analyser) applyCompiler(s status.Status, cf status.Flag, cstr string) {
	if !cf.MatchesStatus(s) {
		return
	}
	if _, ok := a.analysis.Compilers[cstr]; !ok {
		return
	}
	a.analysis.Compilers[cstr].Counts[s]++
}

func (a *analyser) applyTimes(r classification) {
	for cstr, ts := range r.ctimes {
		a.compilerTimes[cstr] = append(a.compilerTimes[cstr], ts...)
	}
	for cstr, ts := range r.rtimes {
		a.runTimes[cstr] = append(a.runTimes[cstr], ts...)
	}
}
