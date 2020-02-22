package planner

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// SubjectProber is the interface of types that allow filling in of subject information.
type SubjectProber interface {
	// ProbeSubject populates subject with information gleaned from investigating its litmus file.
	ProbeSubject(ctx context.Context, subject *subject.Subject) error
}

func (p *Planner) planCorpus(ctx context.Context, seed int64) (subject.Corpus, error) {
	probed, err := p.ProbeCorpus(ctx)
	if err != nil {
		return subject.Corpus{}, err
	}

	// TODO(@MattWindsor91): perform corpus pruning
	return probed.Sample(seed, p.CorpusSize)
}

// ProbeCorpus probes each subject in this planner's corpus file list, producing a Corpus proper.
func (p *Planner) ProbeCorpus(ctx context.Context) (subject.Corpus, error) {
	corpus := subject.NewCorpus(p.InFiles...)

	for i := range corpus {
		if err := p.Source.ProbeSubject(ctx, &corpus[i]); err != nil {
			return corpus, err
		}
	}

	return corpus, nil
}
