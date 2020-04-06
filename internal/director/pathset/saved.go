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
	planBasename       = "plan.toml"
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
	Dirs [subject.NumStatus]string
}

// NewSaved creates a save pathset rooted at root.
func NewSaved(root string) *Saved {
	return &Saved{
		Dirs: [subject.NumStatus]string{
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
	return s.Dirs[subject.FirstBadStatus:]
}

// SubjectDir tries to get the directory for saved subjects for status st and iteration time iterTime.
func (s *Saved) SubjectDir(st subject.Status, iterTime time.Time) (string, error) {
	if st < subject.FirstBadStatus || subject.NumStatus <= st {
		return "", fmt.Errorf("%w: not an 'interesting' status", subject.ErrBadStatus)
	}
	return filepath.Join(
		s.Dirs[st],
		strconv.Itoa(iterTime.Year()),
		strconv.Itoa(int(iterTime.Month())),
		strconv.Itoa(iterTime.Day()),
		iterTime.Format("150405"),
	), nil
}

// PlanFile gets the path to which a final plan file for the test at time iterTime, failing with final status st,
// should be saved.
func (s *Saved) PlanFile(st subject.Status, iterTime time.Time) (string, error) {
	return s.subjectFile(planBasename, st, iterTime)
}

// SubjectTarFile gets the path to which a tarball for compile-failed subject sname,
// from the test at time iterTime and with final status st, should be saved.
func (s *Saved) SubjectTarFile(sname string, st subject.Status, iterTime time.Time) (string, error) {
	return s.subjectFile(sname+tarSuffix, st, iterTime)
}

func (s *Saved) subjectFile(fname string, st subject.Status, iterTime time.Time) (string, error) {
	root, err := s.SubjectDir(st, iterTime)
	if err != nil {
		return "", err
	}
	return filepath.Join(root, fname), nil
}
