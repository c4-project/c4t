// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pathset

import (
	"path/filepath"
	"strconv"
	"time"
)

const (
	segCompileFailures = "compile_fail"
	segFlagged         = "flagged"
	segRunFailures     = "run_fail"
	segTimeouts        = "timeout"
	tarSuffix          = ".tar.gz"
)

// Saved contains the pre-computed paths for saving 'interesting' run results.
type Saved struct {
	// DirCompileFailures stores subjects that ran into compiler failures.
	DirCompileFailures string
	// DirFlagged stores subjects that have been flagged as possibly buggy.
	DirFlagged string
	// DirRunFailures stores subjects that met run failures.
	DirRunFailures string
	// DirTimeouts stores subjects that timed out.
	DirTimeouts string
}

// NewSaved creates a save pathset rooted at root.
func NewSaved(root string) *Saved {
	return &Saved{
		DirCompileFailures: filepath.Join(root, segCompileFailures),
		DirFlagged:         filepath.Join(root, segFlagged),
		DirRunFailures:     filepath.Join(root, segRunFailures),
		DirTimeouts:        filepath.Join(root, segTimeouts),
	}
}

// CompileFailureTarFile gets the path to which a tarball for compile-failed subject sname,
// from the test at time iterTime, should be saved.
func (s *Saved) CompileFailureTarFile(sname string, iterTime time.Time) string {
	return tarFile(s.DirCompileFailures, sname, iterTime)
}

// Flagged gets the path to which a tarball for flagged subject sname,
// from the test at time iterTime, should be saved.
func (s *Saved) FlaggedTarFile(sname string, iterTime time.Time) string {
	return tarFile(s.DirFlagged, sname, iterTime)
}

// RunFailureTarFile gets the path to which a tarball for run-failed subject sname,
// from the test at time iterTime, should be saved.
func (s *Saved) RunFailureTarFile(sname string, iterTime time.Time) string {
	return tarFile(s.DirRunFailures, sname, iterTime)
}

// TimeoutTarFile gets the path to which a tarball for compile-failed subject sname,
// from the test at time iterTime, should be saved.
func (s *Saved) TimeoutTarFile(sname string, iterTime time.Time) string {
	return tarFile(s.DirTimeouts, sname, iterTime)
}

// tarFile decides what to call the tarball of the subject with name sname, from the test at time iterTime,
// relative to root root.
func tarFile(root, sname string, iterTime time.Time) string {
	file := sname + tarSuffix
	return filepath.Join(
		root,
		strconv.Itoa(iterTime.Year()),
		strconv.Itoa(int(iterTime.Month())),
		strconv.Itoa(iterTime.Day()),
		iterTime.Format("150405"),
		file,
	)
}
