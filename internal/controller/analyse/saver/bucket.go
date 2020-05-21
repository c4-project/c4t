// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package saver

import (
	"time"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/plan"
	"github.com/MattWindsor91/act-tester/internal/model/status"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// bucketSaver handles the setup of per-status buckets in an analysis save.
type bucketSaver struct {
	s        status.Status
	plan     *plan.Plan
	paths    *RunPathset
	parent   *Saver
	creation time.Time
}

func (b *bucketSaver) save(c corpus.Corpus) error {
	if err := b.paths.Prepare(); err != nil {
		return err
	}
	if err := b.writePlan(); err != nil {
		return err
	}
	return b.archiveSubjects(c)
}

func (b *bucketSaver) writePlan() error {
	return b.plan.WriteFile(b.paths.FilePlan)
}

func (b *bucketSaver) archiveSubjects(corp corpus.Corpus) error {
	for name, sub := range corp {
		if err := b.archiveSubject(name, sub); err != nil {
			return err
		}
	}
	return nil
}

func (b *bucketSaver) archiveSubject(name string, sub subject.Subject) error {
	st, err := b.makeArchiver(name, sub)
	if err != nil {
		return err
	}
	aerr := st.archive()
	cerr := st.archiver.Close()
	return iohelp.FirstError(aerr, cerr)
}

func (b *bucketSaver) makeArchiver(name string, sub subject.Subject) (*subjectArchiver, error) {
	path := b.paths.SubjectTarFile(name)
	ar, err := b.parent.archiveMaker(path)
	if err != nil {
		return nil, err
	}
	st := subjectArchiver{
		sub:       sub.AddName(name),
		observers: b.parent.observers,
		archiver:  ar,
		path:      path,
	}
	return &st, nil
}
