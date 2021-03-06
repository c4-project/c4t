// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package perturber

import (
	"math/rand"

	"github.com/c4-project/c4t/internal/subject/corpus"
	"github.com/c4-project/c4t/internal/subject/corpus/builder"

	"github.com/c4-project/c4t/internal/plan"
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
	// TODO(@MattWindsor91): the fact that we're reusing the builder observations here is sus.
	obs := lowerToBuilder(p.observers)
	builder.OnBuild(builder.StartMessage(builder.Manifest{Name: "sampled", NReqs: len(c)}), obs...)
	for i, n := range c.Names() {
		s := c[n]
		builder.OnBuild(builder.StepMessage(i, builder.AddRequest(s.AddName(n))), obs...)
	}
	builder.OnBuild(builder.EndMessage(), obs...)
}
