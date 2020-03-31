// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"context"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// SubjectProber is the interface of types that allow filling in of subject information.
type SubjectProber interface {
	// ProbeSubject probes the Litmus test at file litmus, producing a named subject record.
	ProbeSubject(ctx context.Context, litmus string) (subject.Named, error)
}

func (p *Planner) planCorpus(ctx context.Context, rng *rand.Rand, fs []string) (corpus.Corpus, error) {
	c := CorpusPlanner{
		Files:     fs,
		Prober:    p.Source.SProbe,
		Observers: p.Observers,
		Rng:       rng,
		Size:      p.CorpusSize,
	}
	return c.Plan(ctx)
}

// CorpusPlanner contains the state required to plan the corpus part of an initial plan file.
type CorpusPlanner struct {
	// Files contains the files that are to be included in the plan.
	Files []string
	// Observers observe the process of building the corpus.
	Observers []builder.Observer
	// Prober tells the planner how to probe corpus files for specific information.
	Prober SubjectProber
	// Rng is the random number generator to use in corpus sampling.
	Rng *rand.Rand
	// Size is the target size of the corpus.
	Size int
}

func (p *CorpusPlanner) Plan(ctx context.Context) (corpus.Corpus, error) {
	probed, err := p.probe(ctx)
	if err != nil {
		return corpus.Corpus{}, err
	}

	// TODO(@MattWindsor91): perform corpus pruning
	return p.sample(probed)
}

// probe probes each subject in this planner's corpus file list, producing a Corpus proper.
func (p *CorpusPlanner) probe(ctx context.Context) (corpus.Corpus, error) {
	var c corpus.Corpus

	b, berr := p.makeBuilder()
	if berr != nil {
		return nil, berr
	}

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return p.probeInner(ectx, b.SendCh)
	})
	eg.Go(func() error {
		var err error
		c, err = b.Run(ectx)
		return err
	})
	err := eg.Wait()

	return c, err
}

func (p *CorpusPlanner) makeBuilder() (*builder.Builder, error) {
	bc := builder.Config{
		Init:      nil,
		Observers: p.Observers,
		Manifest: builder.Manifest{
			Name:  "plan",
			NReqs: len(p.Files),
		},
	}
	return builder.NewBuilder(bc)
}

func (p *CorpusPlanner) probeInner(ctx context.Context, ch chan<- builder.Request) error {
	for _, f := range p.Files {
		if err := p.probeSubject(ctx, f, ch); err != nil {
			return err
		}
	}
	return nil
}

func (p *CorpusPlanner) probeSubject(ctx context.Context, f string, ch chan<- builder.Request) error {
	s, err := p.Prober.ProbeSubject(ctx, f)
	if err != nil {
		return err
	}
	return builder.AddRequest(&s).SendTo(ctx, ch)
}

func (p *CorpusPlanner) sample(c corpus.Corpus) (corpus.Corpus, error) {
	return c.Sample(p.Rng, p.Size)
}
