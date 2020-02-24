// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// SubjectProber is the interface of types that allow filling in of subject information.
type SubjectProber interface {
	// ProbeSubject probes the Litmus test at file litmus, producing a named subject record.
	ProbeSubject(ctx context.Context, litmus string) (subject.Named, error)
}

func (p *Planner) planCorpus(ctx context.Context, seed int64) (corpus.Corpus, error) {
	probed, err := p.ProbeCorpus(ctx)
	if err != nil {
		return corpus.Corpus{}, err
	}

	// TODO(@MattWindsor91): perform corpus pruning
	return probed.Sample(seed, p.CorpusSize)
}

// ProbeCorpus probes each subject in this planner's corpus file list, producing a Corpus proper.
func (p *Planner) ProbeCorpus(ctx context.Context) (corpus.Corpus, error) {
	c := make(corpus.Corpus, len(p.InFiles))

	for _, f := range p.InFiles {
		s, err := p.Source.ProbeSubject(ctx, f)
		if err != nil {
			return nil, err
		}
		if err := c.Add(s); err != nil {
			return nil, err
		}
	}

	return c, nil
}
