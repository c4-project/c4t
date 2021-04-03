// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package litmus contains model structs and associated functions for Litmus test entries.
package litmus

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/c4-project/c4t/internal/id"

	"github.com/1set/gut/ystring"
)

var (
	// ErrStatDumperNil occurs when we try to use a nil stat dumper with PopulateStatsFrom.
	ErrStatDumperNil = errors.New("stat dumper is nil")
)

// Litmus contains information about a single litmus test file.
type Litmus struct {
	// Path contains the slashpath to the file.
	Path string `json:"path"`

	// Arch contains the architecture of the Litmus test, if it is not a C test.
	Arch id.ID `json:"arch,omitempty"`

	// Stats contains, if available, the statistics set for this litmus test.
	Stats *Statset `json:"stats,omitempty"`
}

// New constructs a litmus record for slashpath path and options os.
//
// Remember to call filepath.FromSlash if needed.
func New(path string, os ...Option) (*Litmus, error) {
	l := Litmus{Path: path, Stats: &Statset{}}
	err := Options(os...)(&l)
	return &l, err
}

// NewOrPanic is like New, but panics if there's an error.
//
// Use in tests only.
func NewOrPanic(path string, os ...Option) *Litmus {
	l, err := New(path, os...)
	if err != nil {
		panic(err)
	}
	return l
}

// HasPath checks if this litmus file actually has a path given.
func (l *Litmus) HasPath() bool {
	return ystring.IsNotBlank(l.Path)
}

// Filepath gets the OS path for this litmus file.
func (l *Litmus) Filepath() string {
	return filepath.Clean(l.Path)
}

// PopulateStats uses s to populate the statistics for this litmus file,
func (l *Litmus) PopulateStats(ctx context.Context, s StatDumper) error {
	if s == nil {
		return ErrStatDumperNil
	}

	l.Stats = &Statset{}
	// TODO(@MattWindsor91): ideally it should be the stat dumper that decides whether to do this or not.
	if !l.IsC() {
		return nil
	}
	return s.DumpStats(ctx, l.Stats, l.Filepath())
}

// IsC checks whether a litmus test targets the C language.
func (l *Litmus) IsC() bool {
	return l.Arch.HasPrefix(id.ArchC)
}
