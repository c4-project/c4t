// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package saver

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse/observer"
	"github.com/MattWindsor91/act-tester/internal/model/normaliser"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Archiver is the interface of types that can archive subject files.
type Archiver interface {
	// ArchiveFile archives the file at rpath, storing it in the archive as wpath and with permissions mode.
	ArchiveFile(rpath, wpath string, mode int64) error

	// An Archiver can be closed; this should free all resources attached to the archiver.
	io.Closer
}

type subjectArchiver struct {
	sub       *subject.Named
	path      string
	observers []observer.Observer
	archiver  Archiver
}

func (s *subjectArchiver) saving() observer.Saving {
	return observer.Saving{
		SubjectName: s.sub.Name,
		Dest:        s.path,
	}
}

func (s *subjectArchiver) archive() error {
	fs, err := filesToArchive(s.sub.Subject)
	if err != nil {
		return err
	}
	for wpath, norm := range fs {
		rpath := norm.Original
		if err := s.rescueNotExistError(s.archiver.ArchiveFile(rpath, wpath, 0744), rpath); err != nil {
			return fmt.Errorf("archiving %q: %w", rpath, err)
		}
	}

	observer.OnSave(s.saving(), s.observers...)
	return nil
}

func filesToArchive(s subject.Subject) (map[string]normaliser.Entry, error) {
	n := normaliser.New("")
	if _, err := n.Normalise(s); err != nil {
		return nil, err
	}
	return n.Mappings, nil
}

func (s *subjectArchiver) rescueNotExistError(err error, rpath string) error {
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	observer.OnSaveFileMissing(s.saving(), rpath, s.observers...)
	return nil
}
