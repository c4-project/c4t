// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package normaliser

import (
	"fmt"
	"path"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
)

// Corpus is a corpus-level normaliser.
type Corpus struct {
	root string

	// Mappings contains the merged mapping table across all subjects.
	Mappings Map

	// BySubject contains the normalisers for each subject, which, in turn, contain the mappings.
	BySubject map[string]*Normaliser
}

// NewCorpus constructs a new corpus normaliser relative to root.
func NewCorpus(root string) *Corpus {
	return &Corpus{
		root:      root,
		Mappings:  make(map[string]Entry),
		BySubject: make(map[string]*Normaliser),
	}
}

// Normalise normalises mappings for each subject in c.
func (n *Corpus) Normalise(c corpus.Corpus) (corpus.Corpus, error) {
	c2 := make(corpus.Corpus, len(c))
	for name, s := range c {
		snorm := New(path.Join(n.root, name))
		ns, err := snorm.Normalise(s)
		if err != nil {
			return nil, fmt.Errorf("normalising %s: %w", name, err)
		}

		c2[name] = *ns
		n.add(name, snorm)
	}
	return c2, nil
}

func (n *Corpus) add(name string, snorm *Normaliser) {
	n.BySubject[name] = snorm
	// Assuming we've constructed the mapping such that there can't be any overlaps.
	for k, v := range snorm.Mappings {
		n.Mappings[k] = v
	}
}
