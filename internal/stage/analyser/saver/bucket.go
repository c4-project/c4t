// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package saver

import (
	"time"

	"github.com/MattWindsor91/c4t/internal/helper/errhelp"

	"github.com/MattWindsor91/c4t/internal/plan"
	"github.com/MattWindsor91/c4t/internal/subject/corpus"
	"github.com/MattWindsor91/c4t/internal/subject/status"
)

// bucketSaver handles the setup of per-status buckets in an analysis save.
type bucketSaver struct {
	parent   *Saver
	s        status.Status
	plan     *plan.Plan
	paths    *RunPathset
	creation time.Time
}

func (b *bucketSaver) save(c corpus.Corpus) error {
	if err := b.paths.Prepare(); err != nil {
		return err
	}
	if err := b.writePlan(c); err != nil {
		return err
	}
	return b.archiveSubjects(c)
}

func (b *bucketSaver) writePlan(c corpus.Corpus) error {
	// TODO(@MattWindsor91): we do this because c has non-normalised names, but it still seems a bit extraneous.
	pn := *b.plan
	pn.Corpus = pn.Corpus.FilterToNames(c.Names()...)
	return pn.WriteFile(b.paths.FilePlan, plan.WriteCompress|plan.WriteHuman)
}

func (b *bucketSaver) archiveSubjects(corp corpus.Corpus) error {
	for name := range corp {
		if err := b.archiveSubject(name); err != nil {
			return err
		}
	}
	return nil
}

func (b *bucketSaver) archiveSubject(name string) error {
	path := b.paths.SubjectTarFile(name)
	ar, err := b.parent.archiveMaker(path)
	if err != nil {
		return err
	}
	nameMap := b.parent.norm.BySubject[name].Mappings
	aerr := ArchiveSubject(ar, name, path, nameMap, b.parent.observers...)
	cerr := ar.Close()
	return errhelp.FirstError(aerr, cerr)
}
