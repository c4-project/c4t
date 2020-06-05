// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package fuzzer contains a test-plan batch fuzzer.
// It relies on the existence of a single-file fuzzer such as act-fuzz.
package fuzzer

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// DefaultSubjectCycles is the default number of fuzz cycles to run per subject.
const DefaultSubjectCycles = 10

// Fuzzer holds the configuration required to fuzz a plan file.
type Fuzzer struct {
	l *log.Logger

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
	if err := checkPlan(p); err != nil {
		return nil, err
	}

	f := Fuzzer{plan: *p, l: iohelp.EnsureLog(c.Logger), conf: *c}

	err := f.checkCount()
	return &f, err
}

func checkConfig(c *Config) error {
	if c == nil {
		return ErrConfigNil
	}
	return c.Check()
}

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	return p.Check()
}

func (f *Fuzzer) checkCount() error {
	nsubjects, nruns := f.count()
	if nsubjects <= 0 {
		return corpus.ErrNone
	}

	// Note that this inequality 'does the right thing' when f.CorpusSize = 0, ie no corpus size requirement.
	csize := f.conf.Quantities.CorpusSize
	if nruns < csize {
		return fmt.Errorf("%w: projected corpus size %d, want %d", corpus.ErrSmall, nruns, csize)
	}

	return nil
}

// Fuzz runs the fuzzer with context ctx.
func (f *Fuzzer) Fuzz(ctx context.Context) (*plan.Plan, error) {
	f.l.Println("preparing directories")
	if err := f.conf.Paths.Prepare(); err != nil {
		return nil, err
	}

	f.l.Println("now fuzzing")
	rng := f.plan.Metadata.Rand()
	fcs, ferr := f.fuzzInner(ctx, rng)
	if ferr != nil {
		return nil, ferr
	}

	return f.sampleAndUpdatePlan(fcs, rng)
}

// sampleAndUpdatePlan samples fcs and places the result in the fuzzer's plan.
func (f *Fuzzer) sampleAndUpdatePlan(fcs corpus.Corpus, rng *rand.Rand) (*plan.Plan, error) {
	f.l.Println("sampling corpus")
	scs, err := fcs.Sample(rng, f.conf.Quantities.CorpusSize)
	if err != nil {
		return nil, err
	}

	f.l.Println("updating plan")
	f.plan.Corpus = scs
	f.plan.Metadata = *plan.NewHeader(plan.UseDateSeed)
	return &f.plan, nil
}

// count counts the number of subjects and individual fuzz runs to expect from this fuzzer.
func (f *Fuzzer) count() (nsubjects, nruns int) {
	nsubjects = len(f.plan.Corpus)
	nruns = f.conf.Quantities.SubjectCycles * nsubjects
	return nsubjects, nruns
}

// fuzz actually does the business of fuzzing.
func (f *Fuzzer) fuzzInner(ctx context.Context, rng *rand.Rand) (corpus.Corpus, error) {
	_, nfuzzes := f.count()

	f.l.Printf("Fuzzing %d inputs\n", len(f.plan.Corpus))

	seeds := f.corpusSeeds(rng)

	m := builder.Manifest{Name: "fuzz", NReqs: nfuzzes}
	bc := builder.Config{Manifest: m, Observers: f.conf.Observers}
	return builder.ParBuild(ctx, f.conf.Quantities.NWorkers, f.plan.Corpus, bc, func(ctx context.Context, s subject.Named, ch chan<- builder.Request) error {
		return f.makeInstance(s, seeds[s.Name], ch).Fuzz(ctx)
	})
}

// corpusSeeds generates a seed for each subject in the fuzzer's corpus using rng.
func (f *Fuzzer) corpusSeeds(rng *rand.Rand) map[string]int64 {
	seeds := make(map[string]int64)
	for n := range f.plan.Corpus {
		seeds[n] = rng.Int63()
	}
	return seeds
}

func (f *Fuzzer) makeInstance(s subject.Named, seed int64, resCh chan<- builder.Request) *Instance {
	return &Instance{
		Driver:        f.conf.Driver,
		StatDumper:    f.conf.StatDumper,
		Subject:       s,
		SubjectCycles: f.conf.Quantities.SubjectCycles,
		Pathset:       f.conf.Paths,
		Rng:           rand.New(rand.NewSource(seed)),
		ResCh:         resCh,
	}
}
