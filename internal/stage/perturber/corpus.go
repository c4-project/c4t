// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package perturber

import (
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/plan"
)

func (p *Perturber) sampleCorpus(rng *rand.Rand, pn *plan.Plan) error {
	nc, err := pn.Corpus.Sample(rng, p.quantities.CorpusSize)
	if err != nil {
		return err
	}
	pn.Corpus = nc
	return nil
}
