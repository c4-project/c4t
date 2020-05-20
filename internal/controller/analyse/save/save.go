// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package save

import (
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/1set/gut/yos"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse/observer"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/plan"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Saver contains the state used when saving 'interesting' subjects.
type Saver struct {
	// Observers is the list of observers.
	Observers []observer.Observer

	// Paths contains the pathset used to save subjects for a particular machine.
	Paths *Pathset
}

// Run runs the saving stage over the analysis a.
// It returns p unchanged; this is for signature compatibility with the other director stages.
func (s *Saver) Run(a analysis.Analysis) error {
	if s.Paths == nil {
		return iohelp.ErrPathsetNil
	}

	p := a.Plan
	if p == nil {
		return fmt.Errorf("when saving analysis: %w", plan.ErrNil)
	}
	creation := p.Metadata.Creation

	for st, c := range a.ByStatus {
		if st < status.FirstBad || len(c) == 0 {
			continue
		}
		b := bucketSaver{s: st, plan: p, parent: s, creation: creation}
		if err := b.save(c); err != nil {
			return err
		}
	}
	return nil
}

type bucketSaver struct {
	s        status.Status
	plan     *plan.Plan
	parent   *Saver
	creation time.Time
}

func (b *bucketSaver) save(c corpus.Corpus) error {
	if err := b.prepareDir(); err != nil {
		return err
	}
	if err := b.writePlan(); err != nil {
		return err
	}
	return b.tarSubjects(c)
}

func (b *bucketSaver) prepareDir() error {
	dir, err := b.parent.Paths.SubjectDir(b.s, b.creation)
	if err != nil {
		return err
	}
	return yos.MakeDir(dir)
}

func (b *bucketSaver) writePlan() error {
	path, err := b.parent.Paths.PlanFile(b.s, b.creation)
	if err != nil {
		return err
	}
	return b.plan.WriteFile(path)
}

func (b *bucketSaver) tarSubjects(corp corpus.Corpus) error {
	for name, sub := range corp {
		if err := b.tarSubject(name, sub); err != nil {
			return err
		}
	}
	return nil
}

func (b *bucketSaver) tarSubject(name string, sub subject.Subject) error {
	tarpath, err := b.parent.Paths.SubjectTarFile(name, b.s, b.creation)
	if err != nil {
		return err
	}

	st := subjectTar{
		sub:       sub.AddName(name),
		path:      tarpath,
		observers: b.parent.Observers,
	}
	return st.tar()
}
