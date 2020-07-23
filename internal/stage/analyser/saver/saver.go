// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package saver contains the part of the analyser that uses the analyser to save failing tests.
package saver

import (
	"errors"
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/normaliser"

	"github.com/MattWindsor91/act-tester/internal/plan/analyser"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/observer"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/plan"
)

// Saver contains the state used when saving 'interesting' subjects.
type Saver struct {
	// archiveMaker is a factory function opening an Archiver on an archive at file path.
	archiveMaker func(path string) (Archiver, error)
	// norm is a normaliser used to translate a corpus's paths to the ones used in its archival.
	norm *normaliser.Corpus
	// observers is the list of observers.
	observers []observer.Observer
	// paths contains the pathset used to save subjects for a particular machine.
	paths *Pathset
}

var (
	ErrArchiveMakerNil = errors.New("archive maker function nil")
	ErrObserverNil     = errors.New("observer nil")
)

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
// It returns p unchanged; this is for signature compatibility with the other director stages.
func (s *Saver) Run(a analyser.Analysis) error {
	p := a.Plan
	if p == nil {
		return fmt.Errorf("when saving analyser: %w", plan.ErrNil)
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
