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

// Analyse runs the analyser with context ctx, on plan p and with options opts.
func Analyse(ctx context.Context, p *plan.Plan, opts ...Option) (*Analysis, error) {
	a, err := newAnalyser(p, opts...)
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

	// filters is the set of filters to use when filtering compiler results.
	filters FilterSet
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
	ch := make(chan subjectAnalysis)
	err := a.corpus.Par(ctx, a.nworkers,
		func(ctx context.Context, named subject.Named) error {
			a.analyseAndSend(ctx, named, ch)
			return nil
		},
		func(ctx context.Context) error {
			return a.build(ctx, ch)
		},
	)
	return err
}

// newAnalyser initialises an analyser for plan p, with workers nworkers.
func newAnalyser(p *plan.Plan, opts ...Option) (*analyser, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}

	lc := len(p.Compilers)
	a := analyser{
		analysis:      newAnalysis(p),
		corpus:        p.Corpus,
		compilerTimes: make(map[string][]time.Duration, lc),
		runTimes:      make(map[string][]time.Duration, lc),
	}
	if err := Options(opts...)(&a); err != nil {
		return nil, err
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
		a.analysis.Compilers[cn] = Compiler{Counts: map[status.Status]int{}, Logs: map[string]string{}, Info: c}
		a.compilerTimes[cn] = []time.Duration{}
		a.runTimes[cn] = []time.Duration{}
	}
}

func (a *analyser) analyseAndSend(ctx context.Context, named subject.Named, ch chan<- subjectAnalysis) {
	select {
	case ch <- a.analyseSubject(named):
	case <-ctx.Done():
	}
}

func (a *analyser) build(ctx context.Context, ch <-chan subjectAnalysis) error {
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

func (a *analyser) apply(r subjectAnalysis) {
	a.analysis.Flags |= r.flags
	a.applyCompilers(r)
	a.applyTimes(r)

	for i := status.Ok; i <= status.Last; i++ {
		a.applyByStatus(i, r)
	}
}

func (a *analyser) applyCompilers(r subjectAnalysis) {
	for cstr, cflag := range r.cflags {
		if _, ok := a.analysis.Compilers[cstr]; !ok {
			// Somehow the analysis is mentioning a compiler whose existence we haven't foreseen.
			continue
		}
		a.analysis.Compilers[cstr].Logs[r.sub.Name] = r.clogs[cstr]

		for i := status.Ok; i <= status.Last; i++ {
			a.applyCompilerStatusCount(i, cflag, cstr)
		}
	}
}

func (a *analyser) applyByStatus(s status.Status, r subjectAnalysis) {
	if !r.flags.MatchesStatus(s) {
		return
	}
	if _, ok := a.analysis.ByStatus[s]; !ok {
		a.analysis.ByStatus[s] = make(corpus.Corpus)
	}
	a.analysis.ByStatus[s][r.sub.Name] = r.sub.Subject
}

func (a *analyser) applyCompilerStatusCount(s status.Status, cf status.Flag, cstr string) {
	if !cf.MatchesStatus(s) {
		return
	}
	a.analysis.Compilers[cstr].Counts[s]++
}

func (a *analyser) applyTimes(r subjectAnalysis) {
	for cstr, ts := range r.ctimes {
		a.compilerTimes[cstr] = append(a.compilerTimes[cstr], ts...)
	}
	for cstr, ts := range r.rtimes {
		a.runTimes[cstr] = append(a.runTimes[cstr], ts...)
	}
}
