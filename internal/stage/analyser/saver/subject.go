// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package saver

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/c4-project/c4t/internal/subject/normaliser"
)

// Archiver is the interface of types that can archive subject files.
type Archiver interface {
	// ArchiveFile archives the file at rpath, storing it in the archive as wpath and with permissions mode.
	ArchiveFile(rpath, wpath string, mode int64) error

	// An Archiver can be closed; this should free all resources attached to the archiver.
	io.Closer
}

//go:generate mockery --name=Archiver

type subjectArchiver struct {
	nameMap     normaliser.Map
	sname, dest string
	observers   []Observer
	archiver    Archiver
}

var (
	// ErrArchiverNil occurs when we try to archive a subject with no archiver.
	ErrArchiverNil = errors.New("archiver nil")
)

// ArchiveSubject archives the subject defined by sname and nameMap to dest via ar, announcing progress to obs.
func ArchiveSubject(ar Archiver, sname, dest string, nameMap normaliser.Map, obs ...Observer) error {
	if ar == nil {
		return ErrArchiverNil
	}
	if len(nameMap) == 0 {
		return nil
	}
	s := subjectArchiver{
		nameMap:   nameMap,
		sname:     sname,
		dest:      dest,
		observers: obs,
		archiver:  ar,
	}
	return s.archive()
}

func (s *subjectArchiver) archive() error {
	OnArchiveStart(s.sname, s.dest, len(s.nameMap), s.observers...)

	i := 0
	for wpath, norm := range s.nameMap {
		if err := s.archiveFile(i, wpath, norm); err != nil {
			return fmt.Errorf("archiving %q: %w", norm.Original, err)
		}
		i++
	}

	OnArchiveFinish(s.sname, s.observers...)
	return nil
}

func (s *subjectArchiver) archiveFile(i int, wpath string, norm normaliser.Entry) error {
	perm := norm.Kind.ArchivePerm()
	rpath := norm.Original
	err := s.archiver.ArchiveFile(rpath, wpath, perm)
	if err != nil {
		return s.rescueNotExistError(err, rpath, i)
	}
	OnArchiveFileAdded(s.sname, rpath, i, s.observers...)
	return nil
}

func (s *subjectArchiver) rescueNotExistError(err error, rpath string, i int) error {
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	OnArchiveFileMissing(s.sname, rpath, i, s.observers...)
	return nil
}
