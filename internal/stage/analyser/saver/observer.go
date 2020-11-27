// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package saver

type ArchiveMessageKind uint8

// Observer represents the observer interface for savers.
type Observer interface {
	// OnArchive lets the observer know that an archive action has occurred.
	OnArchive(s ArchiveMessage)
}

//go:generate mockery --name=Observer

const (
	// ArchiveStart denotes the start of an archival run.
	// Dest contains the destination archive; Index contains the number of files to add.
	ArchiveStart ArchiveMessageKind = iota
	// ArchiveFileAdded states that the file with the given index was present and added.
	ArchiveFileAdded
	// ArchiveFileMissing states that the file with the given index was missing and skipped.
	ArchiveFileMissing
	// ArchiveFinish denotes the end of an archival run.
	ArchiveFinish
)

// ArchiveMessage represents an OnArchive message.
type ArchiveMessage struct {
	// Kind contains the kind of archive message being sent.
	Kind ArchiveMessageKind
	// SubjectName is the name of the subject being archived.
	SubjectName string
	// File is the target file of the archive message.
	File string
	// Index is the index of the file, or the number of files being archived, depending on the message kind.
	Index int
}

// OnArchive sends OnArchive to every instance observer in obs.
func OnArchive(s ArchiveMessage, obs ...Observer) {
	for _, o := range obs {
		o.OnArchive(s)
	}
}

// OnArchiveStart sends an archive start message to every observer in obs.
func OnArchiveStart(sname string, dest string, nfiles int, obs ...Observer) {
	OnArchive(ArchiveMessage{
		Kind:        ArchiveStart,
		SubjectName: sname,
		File:        dest,
		Index:       nfiles,
	}, obs...)
}

// OnArchiveFinish sends an archive finish message to every observer in obs.
func OnArchiveFinish(sname string, obs ...Observer) {
	OnArchive(ArchiveMessage{
		Kind:        ArchiveFinish,
		SubjectName: sname,
	}, obs...)
}

// OnArchiveFileAdded sends an archive 'file added' message to every observer in obs.
// This message includes the subject name sname, the added file, and its position i in the archive.
func OnArchiveFileAdded(sname, file string, i int, obs ...Observer) {
	OnArchive(ArchiveMessage{
		Kind:        ArchiveFileAdded,
		SubjectName: sname,
		File:        file,
		Index:       i,
	}, obs...)
}

// OnArchiveFileMissing sends an archive 'file missing' message to every observer in obs.
// This message includes the subject name sname, the missing file missing, and its position i in the archive.
func OnArchiveFileMissing(sname, missing string, i int, obs ...Observer) {
	OnArchive(ArchiveMessage{
		Kind:        ArchiveFileMissing,
		SubjectName: sname,
		File:        missing,
		Index:       i,
	}, obs...)
}
