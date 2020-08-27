// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package perturber

import (
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus"
	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/plan"
)

func (p *Perturber) sampleCorpus(rng *rand.Rand, pn *plan.Plan) error {
	nc, err := pn.Corpus.Sample(rng, p.quantities.CorpusSize)
	if err != nil {
		return err
	}
	pn.Corpus = nc

	p.announceCorpus(nc)

	return nil
}

func (p *Perturber) announceCorpus(c corpus.Corpus) {
	obs := lowerToBuilder(p.observers)
	builder.OnBuildStart(builder.Manifest{Name: "sampled", NReqs: len(c)}, obs...)
	for i, n := range c.Names() {
		s := c[n]
		builder.OnBuildRequest(i, builder.AddRequest(s.AddName(n)), obs...)
	}
	builder.OnBuildFinish(obs...)
}
