// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package saver

import (
	"errors"
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse/observer"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Saver contains the state used when saving 'interesting' subjects.
type Saver struct {
	// archiveMaker is a factory function opening an Archiver on an archive at file path.
	archiveMaker func(path string) (Archiver, error)
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
	}
	if s.paths == nil {
		return nil, iohelp.ErrPathsetNil
	}
	if s.archiveMaker == nil {
		return nil, ErrArchiveMakerNil
	}
	for _, o := range ops {
		if err := o(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// Option is the type of options to New.
type Option func(*Saver) error

// WithObservers appends obs to the observer list for this saver.
func WithObservers(obs ...observer.Observer) Option {
	return func(s *Saver) error {
		for _, o := range obs {
			if o == nil {
				return ErrObserverNil
			}
			s.observers = append(s.observers, o)
		}
		return nil
	}
}

// Run runs the saving stage over the analysis a.
// It returns p unchanged; this is for signature compatibility with the other director stages.
func (s *Saver) Run(a analysis.Analysis) error {
	p := a.Plan
	if p == nil {
		return fmt.Errorf("when saving analysis: %w", plan.ErrNil)
	}
	creation := p.Metadata.Creation

	for st, c := range a.ByStatus {
		if err := s.runBucket(st, c, p, creation); err != nil {
			return err
		}
	}
	return nil
}

func (s *Saver) runBucket(st status.Status, c corpus.Corpus, p *plan.Plan, creation time.Time) error {
	if st < status.FirstBad || len(c) == 0 {
		return nil
	}
	b := bucketSaver{s: st, plan: p, parent: s, creation: creation}
	return b.save(c)
}
