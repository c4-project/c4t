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
)

// Archiver is the interface of types that can archive subject files.
type Archiver interface {
	// ArchiveFile archives the file at rpath, storing it in the archive as wpath and with permissions mode.
	ArchiveFile(rpath, wpath string, mode int64) error

	// An Archiver can be closed; this should free all resources attached to the archiver.
	io.Closer
}

//go:generate mockery -name=Archiver

type subjectArchiver struct {
	nameMap   normaliser.Map
	saving    observer.Saving
	observers []observer.Observer
	archiver  Archiver
}

var (
	// ErrArchiverNil occurs when we try to archive a subject with no archiver.
	ErrArchiverNil = errors.New("archiver nil")
)

// ArchiveSubject archives the subject defined by saving and nameMap to ar, announcing progress to obs.
func ArchiveSubject(ar Archiver, nameMap normaliser.Map, saving observer.Saving, obs ...observer.Observer) error {
	if ar == nil {
		return ErrArchiverNil
	}
	if len(nameMap) == 0 {
		return nil
	}
	s := subjectArchiver{
		nameMap:   nameMap,
		saving:    saving,
		observers: obs,
		archiver:  ar,
	}
	return s.archive()
}

func (s *subjectArchiver) archive() error {
	for wpath, norm := range s.nameMap {
		if err := s.archiveFile(wpath, norm); err != nil {
			return err
		}
	}

	observer.OnSave(s.saving, s.observers...)
	return nil
}

func (s *subjectArchiver) archiveFile(wpath string, norm normaliser.Entry) error {
	perm := norm.Kind.ArchivePerm()
	rpath := norm.Original
	if err := s.rescueNotExistError(s.archiver.ArchiveFile(rpath, wpath, perm), rpath); err != nil {
		return fmt.Errorf("archiving %q: %w", rpath, err)
	}
	return nil
}

func (s *subjectArchiver) rescueNotExistError(err error, rpath string) error {
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	observer.OnSaveFileMissing(s.saving, rpath, s.observers...)
	return nil
}
