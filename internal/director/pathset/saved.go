// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pathset

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

const (
	segFlagged         = "flagged"
	segCompileFailures = "compile_fail"
	segCompileTimeouts = "compile_timeout"
	segRunFailures     = "run_fail"
	segRunTimeouts     = "run_timeout"
	tarSuffix          = ".tar.gz"
)

// Saved contains the pre-computed paths for saving 'interesting' run results.
type Saved struct {
	// Dirs maps 'interesting' statuses to directories.
	Dirs map[subject.Status]string
}

// NewSaved creates a save pathset rooted at root.
func NewSaved(root string) *Saved {
	return &Saved{
		Dirs: map[subject.Status]string{
			subject.StatusFlagged:        filepath.Join(root, segFlagged),
			subject.StatusCompileFail:    filepath.Join(root, segCompileFailures),
			subject.StatusCompileTimeout: filepath.Join(root, segCompileTimeouts),
			subject.StatusRunFail:        filepath.Join(root, segRunFailures),
			subject.StatusRunTimeout:     filepath.Join(root, segRunTimeouts),
		},
	}
}

// DirList gets the list of directories in the save pathset, ordered by subject number.
func (s *Saved) DirList() []string {
	b := subject.FirstBadStatus
	dirs := make([]string, subject.NumStatus-b)
	for i := b; i < subject.NumStatus; i++ {
		dirs[i-b] = s.Dirs[i]
	}
	return dirs
}

// CompileFailureTarFile gets the path to which a tarball for compile-failed subject sname,
// from the test at time iterTime and with final status st, should be saved.
func (s *Saved) SubjectTarFile(sname string, st subject.Status, iterTime time.Time) (string, error) {
	dir, ok := s.Dirs[st]
	if !ok {
		return "", fmt.Errorf("%w: not an 'interesting' status", subject.ErrBadStatus)
	}
	return tarFile(dir, sname, iterTime), nil
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
