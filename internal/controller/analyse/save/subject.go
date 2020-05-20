// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package save

import (
	"errors"
	"fmt"
	"os"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse/observer"
	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/normalise"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

type subjectTar struct {
	sub       *subject.Named
	path      string
	observers []observer.Observer
}

func (s *subjectTar) saving() observer.Saving {
	return observer.Saving{
		SubjectName: s.sub.Name,
		Dest:        s.path,
	}
}

func (s *subjectTar) tar() error {
	tarfile, err := os.Create(s.path)
	if err != nil {
		return fmt.Errorf("create %s: %w", s.path, err)
	}
	tgz := NewTGZWriter(tarfile)
	werr := s.tarToWriter(tgz)
	cerr := tgz.Close()

	observer.OnSave(s.saving(), s.observers...)

	return iohelp.FirstError(werr, cerr)
}

func (s *subjectTar) tarToWriter(tgz *TGZWriter) error {
	fs, err := filesToTar(s.sub.Subject)
	if err != nil {
		return err
	}
	for wpath, norm := range fs {
		rpath := norm.Original
		if err := s.rescueNotExistError(tgz.TarFile(rpath, wpath, 0744), rpath); err != nil {
			return fmt.Errorf("archiving %q: %w", rpath, err)
		}
	}
	return nil
}

func filesToTar(s subject.Subject) (map[string]normalise.Normalisation, error) {
	n := normalise.NewNormaliser("")
	if _, err := n.Subject(s); err != nil {
		return nil, err
	}
	return n.Mappings, nil
}

func (s *subjectTar) rescueNotExistError(err error, rpath string) error {
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	observer.OnSaveFileMissing(s.saving(), rpath, s.observers...)
	return nil
}
