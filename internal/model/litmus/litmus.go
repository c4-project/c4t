// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package litmus contains model structs and associated functions for Litmus test entries.
package litmus

import (
	"context"
	"path/filepath"

	"github.com/1set/gut/ystring"
)

// Litmus contains information about a single litmus test file.
type Litmus struct {
	// Path contains the slashpath to the file.
	Path string `json:"path"`

	// Stats contains, if available, the statistics set for this litmus test.
	Stats *Statset `json:"stats,omitempty"`
}

// New constructs a litmus record for slashpath path and options os.
func New(path string, os ...Option) *Litmus {
	l := Litmus{Path: path, Stats: &Statset{}}
	Options(os...)(&l)
	return &l
}

// NewWithStats uses s to dump statistics for path, and, if successful, returns both as a litmus entry.
func NewWithStats(ctx context.Context, path string, s StatDumper, os ...Option) (*Litmus, error) {
	l := New(path, os...)
	return l, l.PopulateStats(ctx, s)
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
	l.Stats = &Statset{}
	return s.DumpStats(ctx, l.Stats, l.Filepath())
}
