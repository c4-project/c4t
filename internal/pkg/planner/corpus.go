package planner

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// SubjectProber is the interface of types that allow filling in of subject information.
type SubjectProber interface {
	// ProbeSubject populates subject with information gleaned from investigating its litmus file.
	ProbeSubject(subject *model.Subject) error
}

func (p *Planner) planCorpus(seed int64) (model.Corpus, error) {
	probed, err := p.ProbeCorpus()
	if err != nil {
		return model.Corpus{}, err
	}

	// TODO(@MattWindsor91): perform corpus pruning
	return probed.Sample(seed, p.CorpusSize)
}

// ProbeCorpus probes each subject in this planner's corpus file list, producing a Corpus proper.
func (p *Planner) ProbeCorpus() (model.Corpus, error) {
	corpus := model.NewCorpus(p.InFiles...)

	for i := range corpus {
		if err := p.Source.ProbeSubject(&corpus[i]); err != nil {
			return corpus, err
		}
	}

	return corpus, nil
}
