package planner

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

type SubjectProber interface {
}

func (p *Planner) planCorpus(seed int64) (model.Corpus, error) {
	// TODO(@MattWindsor91): perform corpus pruning
	prunedCorpus := model.NewCorpus(p.Corpus...)

	return prunedCorpus.Sample(seed, p.CorpusSize)
}
