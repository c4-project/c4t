// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package saver

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/status"
	"github.com/MattWindsor91/act-tester/internal/plan"
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
	Dirs [status.Last + 1]string
}

// NewPathset creates a save pathset rooted at root.
func NewPathset(root string) *Pathset {
	return &Pathset{
		Dirs: [...]string{
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

func (s *Pathset) SubjectRun(st status.Status, time time.Time) (*RunPathset, error) {
	if !st.IsBad() {
		return nil, fmt.Errorf("%w: not an 'interesting' status", status.ErrBad)
	}
	return s.run(s.Dirs[st], time), nil
}

func (s *Pathset) run(root string, time time.Time) *RunPathset {
	rroot := runRoot(root, time)
	return &RunPathset{
		DirRoot:  rroot,
		FilePlan: filepath.Join(rroot, planBasename+plan.ExtCompress),
	}
}

// RunPathset represents a pathset containing a saved run.
type RunPathset struct {
	// The root directory of the run's pathset.
	DirRoot string

	// The file to which the run's plan should be saved.
	FilePlan string
}

func runRoot(root string, iterTime time.Time) string {
	return filepath.Join(
		root,
		strconv.Itoa(iterTime.Year()),
		strconv.Itoa(int(iterTime.Month())),
		strconv.Itoa(iterTime.Day()),
		iterTime.Format("15_04_05"),
	)
}

// SubjectTarFile gets the path to which a tarball for subject sname should be saved.
func (s *RunPathset) SubjectTarFile(sname string) string {
	return s.subjectFile(sname + tarSuffix)
}

func (s *RunPathset) subjectFile(fname string) string {
	return filepath.Join(s.DirRoot, fname)
}

// Prepare prepares this pathset by making its directories.
func (s *Pathset) Prepare() error {
	return iohelp.Mkdirs(s.DirList()...)
}

// Prepare prepares this pathset by making its directories.
func (s *RunPathset) Prepare() error {
	return iohelp.Mkdirs(s.DirRoot)
}
