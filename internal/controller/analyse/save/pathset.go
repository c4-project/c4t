// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package save

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
	"github.com/MattWindsor91/act-tester/internal/model/status"
)

const (
	planBasename       = "plan"
	segFlagged         = "flagged"
	segCompileFailures = "compile_fail"
	segCompileTimeouts = "compile_timeout"
	segRunFailures     = "run_fail"
	segRunTimeouts     = "run_timeout"
	tarSuffix          = ".tar.gz"
)

// Pathset contains the pre-computed paths for saving 'interesting' run results.
type Pathset struct {
	// Dirs maps 'interesting' statuses to directories.
	Dirs [status.Num]string
}

// NewPathset creates a save pathset rooted at root.
func NewPathset(root string) *Pathset {
	return &Pathset{
		Dirs: [status.Num]string{
			status.Flagged:        filepath.Join(root, segFlagged),
			status.CompileFail:    filepath.Join(root, segCompileFailures),
			status.CompileTimeout: filepath.Join(root, segCompileTimeouts),
			status.RunFail:        filepath.Join(root, segRunFailures),
			status.RunTimeout:     filepath.Join(root, segRunTimeouts),
		},
	}
}

// DirList gets the list of directories in the save pathset, ordered by subject number.
func (s *Pathset) DirList() []string {
	return s.Dirs[status.FirstBad:]
}

// SubjectDir tries to get the directory for saved subjects for status st and iteration time iterTime.
func (s *Pathset) SubjectDir(st status.Status, iterTime time.Time) (string, error) {
	if st < status.FirstBad || status.Num <= st {
		return "", fmt.Errorf("%w: not an 'interesting' status", status.ErrBad)
	}
	return filepath.Join(
		s.Dirs[st],
		strconv.Itoa(iterTime.Year()),
		strconv.Itoa(int(iterTime.Month())),
		strconv.Itoa(iterTime.Day()),
		iterTime.Format("15_04_05"),
	), nil
}

// PlanFile gets the path to which a final plan file for the test at time iterTime, failing with final status st,
// should be saved.
func (s *Pathset) PlanFile(st status.Status, iterTime time.Time) (string, error) {
	return s.subjectFile(planBasename+plan.Ext, st, iterTime)
}

// SubjectTarFile gets the path to which a tarball for compile-failed subject sname,
// from the test at time iterTime and with final status st, should be saved.
func (s *Pathset) SubjectTarFile(sname string, st status.Status, iterTime time.Time) (string, error) {
	return s.subjectFile(sname+tarSuffix, st, iterTime)
}

func (s *Pathset) subjectFile(fname string, st status.Status, iterTime time.Time) (string, error) {
	root, err := s.SubjectDir(st, iterTime)
	if err != nil {
		return "", err
	}
	return filepath.Join(root, fname), nil
}

// Prepare prepares this pathset by making its directories.
func (s *Pathset) Prepare() error {
	return iohelp.Mkdirs(s.DirList()...)
}
