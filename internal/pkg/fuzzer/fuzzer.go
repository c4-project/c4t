// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package fuzzer contains a test-plan batch fuzzer.
// It relies on the existence of a single-file fuzzer such as act-fuzz.
package fuzzer

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"

	"github.com/sirupsen/logrus"
)

// DefaultSubjectCycles is the default number of fuzz cycles to run per subject.
const DefaultSubjectCycles = 10

var (
	// ErrConfigNil occurs when a fuzzer gets constructed using a nil config.
	ErrConfigNil = errors.New("config nil")

	// ErrDriverNil occurs when the fuzzer tries to use the nil pointer as its single-fuzz driver.
	ErrDriverNil = errors.New("driver nil")
)

// Fuzzer holds the configuration required to fuzz a plan file.
type Fuzzer struct {
	// plan is the plan on which this fuzzer is operating.
	plan plan.Plan

	// conf is the configuration used to build this fuzzer.
	conf Config
}

// New constructs a fuzzer with the config c and plan p.
func New(c *Config, p *plan.Plan) (*Fuzzer, error) {
	if err := checkConfig(c); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, plan.ErrNil
	}

	f := Fuzzer{plan: *p, conf: *c}

	err := f.checkCount()
	return &f, err
}

func checkConfig(c *Config) error {
	if c == nil {
		return ErrConfigNil
	}
	if c.Driver == nil {
		return ErrDriverNil
	}
	if c.Paths == nil {
		return iohelp.ErrPathsetNil
	}
	if c.SubjectCycles <= 0 {
		return fmt.Errorf("%w: non-positive subject cycle amount", corpus.ErrSmallCorpus)
	}
	return nil
}

func (f *Fuzzer) checkCount() error {
	nsubjects, nruns := f.count()
	if nsubjects <= 0 {
		return corpus.ErrNoCorpus
	}

	// Note that this inequality 'does the right thing' when f.CorpusSize = 0, ie no corpus size requirement.
	if nruns < f.conf.CorpusSize {
		return fmt.Errorf("%w: projected corpus size %d, want %d", corpus.ErrSmallCorpus, nruns, f.conf.CorpusSize)
	}

	return nil
}

// Fuzz runs the fuzzer with context ctx.
func (f *Fuzzer) Fuzz(ctx context.Context) (*plan.Plan, error) {
	logrus.Infoln("preparing directories")
	if err := f.conf.Paths.Prepare(); err != nil {
		return nil, err
	}

	logrus.Infoln("now fuzzing")
	rng := f.plan.Header.Rand()
	fcs, ferr := f.fuzz(ctx, rng)
	if ferr != nil {
		return nil, ferr
	}

	return f.sampleAndUpdatePlan(fcs, rng)
}

// sampleAndUpdatePlan samples fcs and places the result in the fuzzer's plan.
func (f *Fuzzer) sampleAndUpdatePlan(fcs corpus.Corpus, rng *rand.Rand) (*plan.Plan, error) {
	logrus.Infoln("sampling corpus")
	scs, err := fcs.Sample(rng, f.conf.CorpusSize)
	if err != nil {
		return nil, err
	}

	logrus.Infoln("updating plan")
	f.plan.Corpus = scs
	f.plan.Header = *plan.NewHeader()
	return &f.plan, nil
}

// count counts the number of subjects and individual fuzz runs to expect from this fuzzer.
func (f *Fuzzer) count() (nsubjects, nruns int) {
	nsubjects = len(f.plan.Corpus)
	nruns = f.conf.SubjectCycles * nsubjects
	return nsubjects, nruns
}

// fuzz actually does the business of fuzzing.
func (f *Fuzzer) fuzz(ctx context.Context, rng *rand.Rand) (corpus.Corpus, error) {
	_, nfuzzes := f.count()

	logrus.Infof("Fuzzing %d inputs\n", len(f.plan.Corpus))

	seeds := make(map[string]int64)
	for n := range f.plan.Corpus {
		seeds[n] = rng.Int63()
	}

	bc := corpus.BuilderConfig{
		NReqs: nfuzzes,
		// TODO(@MattWindsor91): decouple this
		Obs: &corpus.PbObserver{},
	}
	b, resCh, berr := corpus.NewBuilder(bc)
	if berr != nil {
		return nil, berr
	}

	var fcs corpus.Corpus

	err := f.plan.ParCorpus(ctx,
		func(ctx context.Context, s subject.Named) error {
			j := f.makeJob(s, seeds[s.Name], resCh)
			return j.Fuzz(ctx)
		},
		func(ctx context.Context) (err error) {
			fcs, err = b.Run(ctx)
			return err
		})
	return fcs, err
}

func (f *Fuzzer) makeJob(s subject.Named, seed int64, resCh chan<- corpus.BuilderReq) *Job {
	return &Job{
		Driver:        f.conf.Driver,
		Subject:       s,
		SubjectCycles: f.conf.SubjectCycles,
		Pathset:       f.conf.Paths,
		Rng:           rand.New(rand.NewSource(seed)),
		ResCh:         resCh,
	}
}
