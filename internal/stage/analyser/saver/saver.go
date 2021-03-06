// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package saver contains the part of the analyser that uses the analyser to save failing tests.
package saver

import (
	"errors"
	"fmt"
	"time"

	"github.com/c4-project/c4t/internal/subject/normaliser"

	"github.com/c4-project/c4t/internal/plan/analysis"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/helper/iohelp"
	"github.com/c4-project/c4t/internal/plan"
	"github.com/c4-project/c4t/internal/subject/corpus"
)

// Saver contains the state used when saving 'interesting' subjects.
type Saver struct {
	// archiveMaker is a factory function opening an Archiver on an archive at file path.
	archiveMaker func(path string) (Archiver, error)
	// norm is a normaliser used to translate a corpus's paths to the ones used in its archival.
	norm *normaliser.Corpus
	// observers is the list of observers.
	observers []Observer
	// paths contains the pathset used to save subjects for a particular machine.
	paths *Pathset
}

// ErrArchiveMakerNil is the error produced when the archive maker supplied to New is nil.
var ErrArchiveMakerNil = errors.New("archive maker function nil")

// New constructs a saver with the pathset paths, archive maker archiveMaker, and options ops.
func New(paths *Pathset, archiveMaker func(path string) (Archiver, error), ops ...Option) (*Saver, error) {
	s := Saver{
		paths:        paths,
		archiveMaker: archiveMaker,
		norm:         normaliser.NewCorpus(""),
	}
	if s.paths == nil {
		return nil, iohelp.ErrPathsetNil
	}
	if s.archiveMaker == nil {
		return nil, ErrArchiveMakerNil
	}
	err := Options(ops...)(&s)
	return &s, err
}

// Run runs the saving stage over the analyser a.
func (s *Saver) Run(a analysis.Analysis) error {
	p := a.Plan
	if p == nil {
		return fmt.Errorf("when saving analysis: %w", plan.ErrNil)
	}
	creation := p.Metadata.Creation

	np, err := s.normalisePlan(p)
	if err != nil {
		return err
	}

	for st, c := range a.ByStatus {
		if err := s.runBucket(st, c, np, creation); err != nil {
			return err
		}
	}
	return nil
}

func (s *Saver) normalisePlan(p *plan.Plan) (*plan.Plan, error) {
	var err error

	np := *p
	np.Corpus, err = s.norm.Normalise(p.Corpus)
	return &np, err
}

func (s *Saver) runBucket(st status.Status, c corpus.Corpus, np *plan.Plan, creation time.Time) error {
	if !st.IsBad() || len(c) == 0 {
		return nil
	}
	paths, err := s.paths.SubjectRun(st, creation)
	if err != nil {
		return err
	}
	b := bucketSaver{
		parent:   s,
		s:        st,
		plan:     np,
		paths:    paths,
		creation: creation,
	}
	return b.save(c)
}
